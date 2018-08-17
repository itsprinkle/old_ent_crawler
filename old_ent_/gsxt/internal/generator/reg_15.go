package generator

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"sync/atomic"
)

// 15位工商注册号生成器
type RegNo15 struct {
	sync.RWMutex
	prefix      string
	types       []string
	failedCount int32
	maxCount    int32
	recvChan    chan string
	closeChan   chan struct{}
}

func NewRegNo15() *RegNo15 {
	return &RegNo15{
		maxCount:  300,
		recvChan:  make(chan string),
		closeChan: make(chan struct{}),
	}
}

func (n *RegNo15) Count(count int) *RegNo15 {
	n.Lock()
	n.maxCount = int32(count)
	n.Unlock()
	return n
}

func (n *RegNo15) Prefix(prefix string) *RegNo15 {
	n.Lock()
	n.prefix = prefix
	n.Unlock()
	return n
}

func (n *RegNo15) Types(types []string) *RegNo15 {
	n.Lock()
	n.types = types
	n.Unlock()
	return n
}

func (n *RegNo15) Do() {
	go func() {
		for _, typ := range n.types {
			atomic.StoreInt32(&n.failedCount, 0)
			length := 8 - len(typ)
			maxValue := int(math.Pow(10, float64(length)))
			for i := 0; i < maxValue; i++ {
				if atomic.LoadInt32(&n.failedCount) >= n.maxCount {
					break
				}
				format := "%s%s%0" + fmt.Sprint(length) + "d"
				regNo := fmt.Sprintf(format, n.prefix, typ, i)
				n.recvChan <- regNo
			}
		}
		close(n.closeChan)
	}()
}

func (n *RegNo15) Done() bool {
	select {
	case <-n.closeChan:
		return true
	default:
	}
	return false
}

func (n *RegNo15) Fail() {
	atomic.AddInt32(&n.failedCount, 1)
}

func (n *RegNo15) Success() {
	atomic.StoreInt32(&n.failedCount, 0)
}

func (n *RegNo15) Next() string {
	select {
	case no := <-n.recvChan:
		return n.Check(no)
	case <-n.closeChan:
	}
	return ""
}

func (n *RegNo15) mod10(value int) int {
	if value == 10 {
		return value
	}
	return value % 10
}

func (n *RegNo15) Check(in string) string {
	if strings.Contains(in, "N") {
		return fmt.Sprintf("%sX", in[:14])
	}
	input := []byte(in)
	value := 10
	for _, c := range input[:14] {
		value = (n.mod10(value+int(c-'0')) * 2) % 11
	}
	num := 11 - value
	if num >= 10 {
		num -= 10
	}
	return fmt.Sprintf("%s%d", input[:14], num)
}
