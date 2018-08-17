package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	flags "github.com/jessevdk/go-flags"

	"gsxt/credit"
	"gsxt/gsxt/creditd"
	"gsxt/gsxt/internal/generator"
	"tools/link"
)

const (
	StateOK       = "OK"
	StateNotFound = "NOT_FOUND"
	StateFail     = "FAIL"
)

var opts struct {
	TCPAddr       string        `short:"t" long:"tcp" description:"TCP监听地址" default:":8200"`
	HTTPAddr      string        `short:"p" long:"http" description:"HTTP监听地址" default:":8201"`
	Retry         bool          `short:"r" long:"retry" description:"是否打开重试"`
	Generator     string        `short:"g" long:"gen" description:"生成器名" default:"reg_15"`
	QueueCount    int32         `short:"q" long:"queue" description:"客户端任务数量" default:"10"`
	IteratorCount int           `short:"i" long:"iterator" description:"reg_15的迭代次数" default:"300"`
	SearchCount   int           `short:"s" long:"search" description:"搜索次数" default:"3"`
	DetailCount   int           `short:"d" long:"detail" description:"获取详情次数" default:"3"`
	Name          string        `short:"n" long:"name" description:"地区名" required:"true"`
	CodeFile      string        `short:"c" long:"code" description:"区划号码文件" required:"true" value-name:"FILE"`
	LoadFile      string        `short:"l" long:"load" description:"已经搜索过文件" value-name:"FILE"`
	Directory     string        `short:"y" long:"dir" description:"保存目录" default:"out"`
	FlushInterval time.Duration `short:"f" long:"flush" description:"刷新时间" default:"5s"`
}

var stats struct {
	StartAt      string `json:"start_at"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	ClientCount  int32  `json:"client_count"`
	TotalCount   int32  `json:"total_count"`
	SuccessCount int32  `json:"success_count"`
	FailCount    int32  `json:"fail_count"`
	RetryCount   int32  `json:"retry_count"`
	TotalSpeed   int32  `json:"total_speed"`
	SuccessSpeed int32  `json:"success_speed"`
}

var (
	loadMap   = make(map[string]bool)
	sendChan  = make(chan creditd.Request)
	writeChan = make(chan creditd.Response)
	retryChan = make(chan creditd.Request, 100)
	stopChan  = make(chan struct{})

	channel = link.NewUint64Channel()
	task    = creditd.NewTask(time.Second * 10)
	gen     generator.Generator
	lock    sync.Mutex

	m2 = map[string]string{
		"1":  "anhui",
		"2":  "beijing",
		"3":  "fujian",
		"4":  "gansu",
		"5":  "guangxi",
		"6":  "hainan",
		"7":  "hebei",
		"8":  "heilongjiang",
		"9":  "henan",
		"10": "hubei",
		"11": "hunan",
		"12": "jiangsu",
		"13": "jilin",
		"14": "liaoning",
		"15": "ningxia",
		"16": "qinghai",
		"17": "shandong",
		"18": "shanghai",
		"19": "shanxi",
		"20": "tianjin",
		"21": "xinjiang",
		"22": "xizang",
		"23": "yunnan",
		"24": "zongju",
		"25": "guangdong",
		"26": "chongqing",
		"27": "zhejiang",
		"28": "sichuan",
		"29": "guizhou",
		"30": "neimenggu",
		"31": "xianxi",
		"32": "jiangxi",
	}
)

// 根据不同的生成器类型分发任务
func sendLoop() {
	if opts.Generator == "reg_15" {
		file, err := os.Open(opts.CodeFile)
		if err != nil {
			log.Fatalf("ERROR: 打开文件%s失败 - %v", opts.CodeFile, err)
		}
		defer file.Close()
		r := bufio.NewReader(file)
		for line, _, err := r.ReadLine(); err == nil; line, _, err = r.ReadLine() {
			if line == nil {
				continue
			}
			prefix := string(line)
			log.Printf("INFO: 开始搜索区划码 - %s", prefix)
			gen = generator.MustGet15(opts.Name, prefix, opts.IteratorCount)
			stats.Code = prefix
			for next := gen.Next(); next != ""; next = gen.Next() {
				if next != "" {
					if val, ok := loadMap[next]; ok {
						if val {
							gen.Success()
						} else {
							gen.Fail()
						}
						continue
					}

					req := creditd.Request{
						Name:        opts.Name,
						Keyword:     next,
						SearchCount: opts.SearchCount,
						DetailCount: opts.DetailCount,
					}
					sendChan <- req
				}
			}
		}
	} else if opts.Generator == "line" {
		gen = generator.MustGetLine(opts.CodeFile)
		for next := gen.Next(); next != ""; next = gen.Next() {
			if next == "" || strings.Contains(next, "not found") {
				continue
			}
			keyword := strings.SplitN(next, "\t", 2)[0]
			if val, ok := loadMap[keyword]; ok {
				if val {
					gen.Success()
				} else {
					gen.Fail()
				}
				continue
			}

			req := creditd.Request{
				Name:        opts.Name,
				Keyword:     next,
				SearchCount: opts.SearchCount,
				DetailCount: opts.DetailCount,
			}
			sendChan <- req
		}
	} else if opts.Generator == "mysql" {
		gen = generator.MustGetMysql()
		for next := gen.Next(); next != ""; next = gen.Next() {
			if next == "" {
				continue
			}
			//关键词、省份、重复次数、id
			spli := strings.SplitN(next, "+", 3)
			keyword := spli[0]
			province := m2[spli[1]]
			req := creditd.Request{
				Extra:       next,
				Name:        province,
				Keyword:     keyword,
				SearchCount: opts.SearchCount,
				DetailCount: opts.DetailCount,
			}
			sendChan <- req
		}
	}
}

// 将响应数据写入硬盘
func writeLoop() {
	dpath := filepath.Join(opts.Directory, "detail", opts.Name+".txt")
	dfile, err := os.OpenFile(dpath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("ERROR: 打开文件%s失败 - %v", dpath, err)
	}
	defer dfile.Close()
	dw := bufio.NewWriterSize(dfile, 10240)

	cpath := filepath.Join(opts.Directory, "code", opts.Name+".txt")
	cfile, err := os.OpenFile(cpath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("ERROR: 打开文件%s失败 - %v", cpath, err)
	}
	defer cfile.Close()
	cw := bufio.NewWriterSize(cfile, 4096)

	go func() {
		select {
		case <-time.After(opts.FlushInterval):
			cw.Flush()
			dw.Flush()
		}
	}()

	for {
		select {
		case <-stopChan:
			cw.Flush()
			dw.Flush()
			return
		case resp := <-writeChan:
			var t []string
			fmt.Println("writeChan>>>>>>>>>>>>>>", resp.Extra, resp.Name, resp.State)
			if opts.Name == "mysql" {
				t = strings.Split(resp.Extra, "+") //keyword,province,crawler_count,id
			}
			if resp.State == StateOK {
				for _, key := range resp.Keys {
					cw.WriteString(fmt.Sprintf("%s\t%s\n", resp.Keyword, key))
				}

				//				for _, info := range resp.Infos {
				//					_, err := json.Marshal(info)
				//					if err == nil {
				//						//creditd.AddMsg(&info, opts.Name)
				//						//dw.WriteString(fmt.Sprintf("%s\n", body))
				//					}
				//				}

				if opts.Name == "mysql" {
					creditd.InsertMsgInfoByV2(resp.Infos, m2[t[1]])
				} else {
					creditd.InsertMsgInfoByV2(resp.Infos, opts.Name)
				}

				if opts.Name == "mysql" {
					creditd.UpdateKeywordStatus(t[3], "3", t[2])
				}
			} else if resp.State == StateNotFound {
				cw.WriteString(fmt.Sprintf("%s\tnot found\n", resp.Keyword))
				if opts.Name == "mysql" {
					creditd.UpdateKeywordStatus(t[3], "2", t[2])
				}
			}
		}
	}
}

// 重试那些超时未返回的数据
func retryLoop() {
	timeoutChan := task.TimeoutChan()
	for {
		select {
		case r := <-timeoutChan:
			req := r.(creditd.Request)
			log.Printf("WARN: 超时没返回数据，重试%s-%s", req.Name, req.Keyword)
			atomic.AddInt32(&stats.RetryCount, 1)
			retryChan <- req
		}
	}
}

// 计算采集速度
// 每5秒钟更新一次
func speedLoop() {
	ticket := time.NewTicker(time.Second * 5)
	successCount := atomic.LoadInt32(&stats.SuccessCount)
	totalCount := atomic.LoadInt32(&stats.TotalCount)

	for {
		select {
		case <-ticket.C:
			atomic.StoreInt32(&stats.ClientCount, int32(channel.Len()))
			atomic.StoreInt32(&stats.TotalSpeed, (stats.TotalCount-totalCount)/5)
			atomic.StoreInt32(&stats.SuccessSpeed, (stats.SuccessCount-successCount)/5)
			successCount = atomic.LoadInt32(&stats.SuccessCount)
			totalCount = atomic.LoadInt32(&stats.TotalCount)
		}
	}
}

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
				if resp.State == StateOK {
					writeChan <- resp
					gen.Success()
					atomic.AddInt32(&stats.SuccessCount, 1)
				} else if resp.State == StateNotFound {
					writeChan <- resp
					gen.Fail()
				} else if resp.State == StateFail {
					atomic.AddInt32(&stats.FailCount, 1)
					// 当访问异常的时候，是否立即重试
					if opts.Retry {
						retryChan <- creditd.Request{
							Name:        resp.Name,
							Keyword:     resp.Keyword,
							SearchCount: opts.SearchCount,
							DetailCount: opts.DetailCount,
						}
					}
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
		case req = <-retryChan:
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
	stats.Name = opts.Name
	server, err := link.Serve("tcp", opts.TCPAddr, link.Bufio(link.GOB()))
	if err != nil {
		log.Fatalf("ERROR: 监听地址%s失败 - %v", opts.TCPAddr, err)
	}
	defer server.Close()
	log.Printf("INFO: 监听地址%s成功", opts.TCPAddr)

	if opts.LoadFile == "" {
		if opts.Generator == "reg_15" {
			loadFile := filepath.Join(opts.Directory, "code", opts.Name+".txt")
			if _, err := os.Stat(loadFile); err == nil {
				opts.LoadFile = loadFile
			}
		}
	}

	if opts.LoadFile != "" {
		file, err := os.Open(opts.LoadFile)
		if err != nil {
			log.Fatalf("ERROR: 打开文件%s失败 - %s", opts.LoadFile, err)
		}
		rr := bufio.NewReader(file)
		for line, _, err := rr.ReadLine(); err == nil; line, _, err = rr.ReadLine() {
			keyword := strings.Split(string(line), "\t")[0]
			if strings.HasPrefix(string(line), "{") {
				var info credit.InfoV2
				println(string(line))
				err = json.NewDecoder(bytes.NewReader(line)).Decode(&info)
				if err != nil {
					log.Printf("WARN: 解析出错 - %v", err)
					continue
				}
				keyword = info.Business.Base.RegNo
			}
			// 去掉重复计算
			if _, ok := loadMap[keyword]; ok {
				continue
			}

			if strings.Contains(string(line), "not found") {
				loadMap[keyword] = false
			} else {
				loadMap[keyword] = true
				atomic.AddInt32(&stats.SuccessCount, 1)
			}
			atomic.AddInt32(&stats.TotalCount, 1)
		}
		log.Printf("INFO: 加载已经搜索的文件, 长度为%d", len(loadMap))
		file.Close()
	}

	// 提供HTTP监控
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		body, _ := json.MarshalIndent(stats, "", "    ")
		rw.Write(body)
	})
	go func() {
		log.Fatalln(http.ListenAndServe(opts.HTTPAddr, nil))
	}()

	stats.StartAt = time.Now().Format("2006-01-02 15:04:05")
	go sendLoop()
	go writeLoop()
	go retryLoop()
	go speedLoop()

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
	close(stopChan)
	log.Printf("INFO: 搜索完毕")
	time.Sleep(time.Second * 2)
}
