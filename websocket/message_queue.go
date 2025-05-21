package websocket

import (
	"sync"
)

const EMessageQueueDefaultCap = 4

type MessageQueue struct {
	list     [][]byte
	mu       sync.Mutex
	listCond *sync.Cond
}

func (queue *MessageQueue) Add(msg []byte) {
	queue.mu.Lock()
	queue.list = append(queue.list, msg)
	queue.mu.Unlock()

	queue.listCond.Signal()
}

func (queue *MessageQueue) Reset() {
	queue.mu.Lock()
	queue.reset()
	queue.mu.Unlock()
}

func (queue *MessageQueue) reset() {
	queue.list = queue.list[0:0]
}

func (queue *MessageQueue) Pick(retList *[][]byte) bool {
	var exit bool

	queue.mu.Lock()
	defer queue.mu.Unlock()

	// 防止通道为0阻塞
	for len(queue.list) == 0 {
		queue.listCond.Wait()
	}

	for _, data := range queue.list {
		if data == nil {
			exit = true
			break
		} else {
			*retList = append(*retList, data)
		}
	}

	queue.reset()

	return exit
}

func NewMessageQueue() *MessageQueue {
	return NewMessageQueueWithCapacity(EMessageQueueDefaultCap)
}

func NewMessageQueueWithCapacity(cap int) *MessageQueue {
	queue := &MessageQueue{}
	queue.list = make([][]byte, 0, cap)
	queue.listCond = sync.NewCond(&queue.mu)
	return queue
}
