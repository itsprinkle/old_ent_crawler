package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	"gsxt/credit"
	"gsxt/gsxt/creditd"
	"tools/link"

	flags "github.com/jessevdk/go-flags"
)

const (
	StateOK       = "OK"
	StateNotFound = "NOT_FOUND"
	StateFail     = "FAIL"
)

var opts struct {
	TCPAddr   string        `short:"t" long:"tcp" description:"调度地址" default:":8200"`
	BreakAddr string        `short:"b" long:"break" description:"验证码地址" default:"139.198.4.104:9801"`
	Interval  time.Duration `short:"i" long:"interval" description:"重启ADSL间隔" default:"2m"`
	Restart   bool          `short:"r" long:"restart" description:"是否固定间隔时间定时重启ADSL"`
	Timeout   time.Duration `short:"o" long:"timeout" description:"获取信息的超时时间" default:"2m"`
}

var (
	pool   = creditd.NewPool(20)
	status int32
	sess   *link.Session
)

func adsl() {
	if !opts.Restart {
		return
	}
	if atomic.CompareAndSwapInt32(&status, 0, 1) {
		var sem = make(chan byte)
		go func() {
			out, err := exec.Command("sh", "-c", "pppoe-stop").Output()
			if err != nil {
				println("pppoe-stop", err.Error())
			} else {
				println("pppoe-stop", string(out))
			}
			out, err = exec.Command("sh", "-c", "pppoe-start").Output()
			if err == nil {
				println("pppoe-start", string(out))
			} else {
				println("pppoe-start", err.Error())
			}
			time.Sleep(time.Second * 1)
			select {
			case sem <- 1:
			default:
			}
		}()
		select {
		case <-sem:
		case <-time.After(time.Second * 15):
		}
		if sess != nil {
			sess.SetReadDeadline(time.Now().Add(time.Second))
		}
		atomic.StoreInt32(&status, 0)
	} else {
		time.Sleep(time.Second * 8)
	}
}

func handle(session *link.Session, req creditd.Request) {
	client := pool.Borrow(req.Name)
	defer func() {
		pool.Return(req.Name, client)
	}()

	client = credit.MustGet(req.Name)

	fmt.Println("req>>>", req.Name, req.Extra)
	var resp creditd.Response
	resp.Extra = req.Extra
	resp.Name = req.Name
	resp.Keyword = req.Keyword
	for i := 0; i < req.SearchCount; i++ {
		log.Printf("INFO: 开始搜索%d - %s", i, req.Keyword)
		var (
			keys  []string
			infos []credit.InfoV2
			err   error
		)
		li := strings.SplitN(req.Keyword, "\t", 2)
		if len(li) == 1 {
			keys, infos, err = creditd.Get(client, req.Keyword, req.SearchCount, req.DetailCount, -1)
		} else {
			var info credit.InfoV2
			info, err = creditd.GetInfo(client, li[0], li[1], req.DetailCount)
			keys = append(keys, li[1])
			infos = append(infos, info)
		}
		if err != nil {
			if err == credit.ErrNotFound {
				resp.State = StateNotFound
				session.Send(resp)
				return
			}
			log.Printf("WARN: 搜索%s:%d失败 - %v", resp.Keyword, i, err)

			if err == credit.ErrOutOfLimit || strings.Contains(err.Error(), "unreachable") {
				adsl()
				continue
			}
			time.Sleep(time.Second * 2)
			continue
		}
		//没有报not Found错误但是查询infos结果未零——强制返回错误
		if len(infos) == 0 {
			resp.State = StateFail
			session.Send(resp)
			return
		}
		// 如果搜索成功
		resp.State = StateOK
		resp.Keys = keys
		resp.Infos = infos
		session.Send(resp)
		return
	}
	resp.State = StateFail
	if session != nil {
		session.Send(resp)
	}
	return
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatalf("ERROR: 解析配置错误 - %s", err)
	}
	credit.DefaultAddr = opts.BreakAddr
	credit.DefaultTimeout = opts.Timeout
	creditd.InfoTimeout = opts.Timeout

	go func() {
		ticker := time.NewTicker(opts.Interval)
		for {
			select {
			case <-ticker.C:
				if opts.Restart {
					log.Printf("INFO: 定时尝试重启adsl")
				}
				adsl()
			}
		}
	}()

login:
	sess, err = link.Connect("tcp", opts.TCPAddr, link.Bufio(link.GOB()))
	if err != nil {
		log.Printf("ERROR: 连接服务端%s失败，10s后重连 - %s", opts.TCPAddr, err)
		time.Sleep(time.Second * 10)
		adsl()
		goto login
	}
	log.Printf("INFO: 连接服务器%s成功", opts.TCPAddr)

	for {
		var req creditd.Request
		//sess.SetReadDeadline(time.Now().Add(time.Minute * 20))
		if err := sess.Receive(&req); err != nil {
			log.Printf("ERROR: 接收消息失败 - %s", err)
			goto login
		}

		go handle(sess, req)
	}
}
