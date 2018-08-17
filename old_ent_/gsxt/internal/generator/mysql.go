package generator

import (
	"gsxt/gsxt/creditd"
	"time"
)

type Mysql struct {
	count     int
	msg       []string
	recvChan  chan string
	closeChan chan struct{}
}

func NewMysql() (*Mysql, error) {
	return &Mysql{
		count:     0,
		msg:       make([]string, 0),
		recvChan:  make(chan string),
		closeChan: make(chan struct{}),
	}, nil
}

func (l *Mysql) Do() {
	go func() {
		for {
			if (len(l.msg) > 0 && l.count >= len(l.msg)) || len(l.msg) == 0 {
				l.msg = creditd.QueryKeywordMsg()
				l.count = 0
				if len(l.msg) == 0 {
					time.Sleep(600)
					continue
				}
			}
			l.recvChan <- l.msg[l.count]
			l.count += 1
		}
	}()
}

func (l *Mysql) Done() bool {
	select {
	case <-l.closeChan:
		return true
	default:
	}
	return false
}

func (l *Mysql) Fail()    {}
func (l *Mysql) Success() {}

func (l *Mysql) Next() string {
	select {
	case no := <-l.recvChan:
		return no
	case <-l.closeChan:
	}
	return ""
}
