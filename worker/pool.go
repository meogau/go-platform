package worker

import (
	"fmt"
	"sync"
	"time"
)

type Pool struct {
	capacity       int
	messageChannel chan string
	wg             sync.WaitGroup
}

func NewWorkerPool(capacity int, messageQueueSize int) *Pool {
	return &Pool{
		capacity:       capacity,
		messageChannel: make(chan string, messageQueueSize),
	}
}

func (pool *Pool) PushMessage(message string) {
	pool.messageChannel <- message
}

func (pool *Pool) Run() {
	for w := 1; w <= pool.capacity; w++ {
		pool.wg.Add(1)
		go pool.handleMessage(w, pool.messageChannel)
	}
}

func (pool *Pool) handleMessage(id int, channel <-chan string) {
	defer pool.wg.Done()
	for message := range channel {
		fmt.Printf("worker %v is processing message %v\n", id, message)
		time.Sleep(time.Second * 10)
	}
}

// Shutdown implement graceful shutdown
func (pool *Pool) Shutdown() {
	fmt.Printf("Closing worker pool\n")
	close(pool.messageChannel)
	pool.wg.Wait()
	fmt.Printf("Closed worker pool\n")
}
