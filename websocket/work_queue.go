package websocket

import (
	"sync"
)

// WorkQueue struct    工作队列
type WorkQueue struct {
	list []func()
	mu   sync.Mutex
}

// Add method    添加工作队列
func (queue *WorkQueue) Add(work func()) {
	queue.mu.Lock()
	queue.list = append(queue.list, work)
	queue.mu.Unlock()
}

// Reset method    重置工作队列
func (queue *WorkQueue) Reset() {
	queue.mu.Lock()
	queue.reset()
	queue.mu.Unlock()
}

func (queue *WorkQueue) reset() {
	queue.list = make([]func(), 0)
}

// Dump method    获取队列列表并销毁
func (queue *WorkQueue) Dump() []func() {
	queue.mu.Lock()
	retList := queue.list
	queue.reset()
	queue.mu.Unlock()
	return retList
}

// NewWorkQueue function    新建工作队列
func NewWorkQueue() *WorkQueue {
	return &WorkQueue{
		list: make([]func(), 0),
	}
}
