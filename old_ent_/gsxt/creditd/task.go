package creditd

import (
	"sync"
	"time"
)

type Item struct {
	Value interface{}
	time  time.Time
}

type Task struct {
	sync.RWMutex
	every       time.Duration
	timeoutChan chan interface{}
	data        map[string]*Item
}

func NewTask(every time.Duration) *Task {
	t := &Task{
		every:       every,
		timeoutChan: make(chan interface{}, 2000),
		data:        make(map[string]*Item),
	}
	go t.do()
	return t
}

func (t *Task) TimeoutChan() <-chan interface{} {
	return t.timeoutChan
}

// 定时检测超时任务
// 如果发现timeoutChan写满, 将数据丢弃
// 使用者应该监听TimeoutChan
func (t *Task) do() {
	ticket := time.NewTicker(t.every)
	for {
		select {
		case <-ticket.C:
			for k, v := range t.copy() {
				if time.Now().Sub(v.time) > 0 {
					select {
					case t.timeoutChan <- v.Value:
					default:
					}
					t.Del(k)
				}
			}
		}
	}
}

func (t *Task) copy() (data map[string]*Item) {
	t.RLock()
	defer t.RUnlock()
	data = make(map[string]*Item)
	for k, v := range t.data {
		data[k] = v
	}
	return data
}

func (t *Task) Set(key string, value interface{}) {
	t.Lock()
	defer t.Unlock()
	t.data[key] = &Item{Value: value, time: time.Now().Add(time.Minute * 5)}
}

func (t *Task) Get(key string) (value interface{}) {
	t.RLock()
	defer t.RUnlock()
	return t.data[key]
}

func (t *Task) Del(key string) {
	t.Lock()
	defer t.Unlock()
	delete(t.data, key)
}
