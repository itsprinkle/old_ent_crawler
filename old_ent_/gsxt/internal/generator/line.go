package generator

import (
	"bufio"
	"errors"
	"log"
	"os"
)

type Line struct {
	filepath  string
	recvChan  chan string
	closeChan chan struct{}
}

func NewLine(filepath string) (*Line, error) {
	info, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, errors.New("只能打开文件")
	}
	return &Line{
		filepath:  filepath,
		recvChan:  make(chan string),
		closeChan: make(chan struct{}),
	}, nil
}

func (l *Line) Do() {
	go func() {
		file, err := os.Open(l.filepath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		r := bufio.NewReader(file)
		for line, _, err := r.ReadLine(); err == nil; line, _, err = r.ReadLine() {
			if len(line) == 0 {
				continue
			}
			l.recvChan <- string(line)
		}
		close(l.closeChan)
	}()
}

func (l *Line) Done() bool {
	select {
	case <-l.closeChan:
		return true
	default:
	}
	return false
}

func (l *Line) Fail()    {}
func (l *Line) Success() {}

func (l *Line) Next() string {
	select {
	case no := <-l.recvChan:
		return no
	case <-l.closeChan:
	}
	return ""
}
