package main

import (
	"encoding/json"
	"gsxt/gsxt/creditd"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"tools/link"

	flags "github.com/jessevdk/go-flags"
)

const (
	StateOK       = "OK"
	StateNotFound = "NOT_FOUND"
	StateFail     = "FAIL"
)

var opts struct {
	TCPAddr     string        `short:"t" long:"tcp" description:"TCP监听地址" default:":8200"`
	HTTPAddr    string        `short:"p" long:"http" description:"HTTP监听地址" default:":8201"`
	SearchCount int           `short:"s" long:"search" description:"搜索次数" default:"2"`
	QueueCount  int32         `short:"q" long:"queue" description:"客户端任务数量" default:"10"`
	DetailCount int           `short:"d" long:"detail" description:"获取详情次数" default:"1"`
	Timeout     time.Duration `short:"o" long:"timeout" description:"超时事件" default:"1m"`
}

var stats struct {
	StartAt      string `json:"start_at"`
	ClientCount  int32  `json:"client_count"`
	TotalCount   int32  `json:"total_count"`
	SuccessCount int32  `json:"success_count"`
	FailCount    int32  `json:"fail_count"`
	TotalSpeed   int32  `json:"total_speed"`
	SuccessSpeed int32  `json:"success_speed"`
}

type token struct {
	resp  creditd.Response
	close chan struct{}
}

func newToken() *token {
	return &token{
		close: make(chan struct{}),
	}
}

var (
	tokenMap = make(map[string]*token)
	sendChan = make(chan creditd.Request)
	exitChan = make(chan struct{})

	channel = link.NewUint64Channel()
	task    = creditd.NewTask(time.Second * 10)
	lock    sync.Mutex
)

func closeChan(exitChan chan struct{}) {
	lock.Lock()
	defer lock.Unlock()

	select {
	case <-exitChan:
	default:
		close(exitChan)
	}
}

func handle(session *link.Session) {
	exitChan := make(chan struct{})
	var queueCount int32

	go func(session *link.Session) {
		for {
			select {
			case <-exitChan:
				return
			default:
			}

			var resp creditd.Response
			session.SetReadDeadline(time.Now().Add(time.Second * 180))
			if err := session.Receive(&resp); err != nil {
				closeChan(exitChan)
				return
			}
			log.Printf("INFO: 客户端%s响应%s-%s-%s-%d", session.RemoteAddr(), resp.Name, resp.Keyword, resp.State, len(resp.Keys))

			select {
			case <-exitChan:
				return
			default:
			}

			atomic.AddInt32(&queueCount, -1)
			if task.Get(resp.Keyword) != nil {
				atomic.AddInt32(&stats.TotalCount, 1)
				token := tokenMap[resp.Keyword]
				if token != nil {
					token.resp = resp
					close(token.close)
				}
				task.Del(resp.Keyword)
			}
		}
	}(session)

	for {
		// 如果客户端队列消息已经满
		for atomic.LoadInt32(&queueCount) >= opts.QueueCount {
			time.Sleep(time.Millisecond * 100)
		}
		atomic.AddInt32(&queueCount, 1)

		var req creditd.Request
		select {
		case req = <-sendChan:
		case <-exitChan:
			return
		}

		log.Printf("INFO: 客户端%s分发任务%s-%s", session.RemoteAddr(), req.Name, req.Keyword)
		if err := session.Send(req); err != nil {
			closeChan(exitChan)
			break
		}
		task.Set(req.Keyword, req)
	}
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatalf("ERROR: 解析配置错误 - %s", err)
	}

	server, err := link.Serve("tcp", opts.TCPAddr, link.Bufio(link.GOB()))
	if err != nil {
		log.Fatalf("ERROR: 监听地址%s失败 - %v", opts.TCPAddr, err)
	}
	defer server.Close()
	log.Printf("INFO: 监听地址%s成功", opts.TCPAddr)

	http.HandleFunc("/search", func(rw http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		name := query.Get("name")
		keyword := query.Get("keyword")

		if name == "" || keyword == "" {
			rw.WriteHeader(400)
			rw.Write([]byte("ERR_INVALID_ARGS"))
			return
		}

		token := tokenMap[keyword]
		if token == nil {
			req := creditd.Request{
				Name:        name,
				Keyword:     keyword,
				SearchCount: opts.SearchCount,
				DetailCount: opts.DetailCount,
			}
			sendChan <- req
			token = newToken()
			tokenMap[keyword] = token
		}

		select {
		case <-token.close:
			if token.resp.State == StateOK {
				atomic.AddInt32(&stats.SuccessCount, 1)
			} else if token.resp.State == StateNotFound {
			} else if token.resp.State == StateFail {
				atomic.AddInt32(&stats.FailCount, 1)
			}
			body, _ := json.Marshal(token.resp)
			rw.Write(body)
		case <-time.After(opts.Timeout):
			rw.WriteHeader(400)
			rw.Write([]byte("ERR_CLIENT_TIMEOUT"))
		}
		delete(tokenMap, keyword)
		return
	})
	go func() {
		log.Fatalln(http.ListenAndServe(opts.HTTPAddr, nil))
	}()

	go func() {
		for {
			session, err := server.Accept()
			if err != nil {
				log.Printf("ERROR: 客户端连接失败 - %s", err)
				continue
			}

			channel.Put(session.Id(), session)
			log.Printf("INFO: 客户端%s已经连接, 当前客户端数量为%d", session.RemoteAddr(), channel.Len())

			go handle(session)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	log.Printf("INFO: 服务退出...")
	time.Sleep(time.Second * 1)
}
