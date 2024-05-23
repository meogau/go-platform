package worker

import (
	"fmt"
	"sync"
)

type Pool struct {
	capacity       int
	messageChannel chan interface{}
	wg             sync.WaitGroup
	processFunc    func(message interface{}) error
}

func NewWorkerPool(capacity int, messageQueueSize int, processFunc func(message interface{}) error) *Pool {
	return &Pool{
		capacity:       capacity,
		messageChannel: make(chan interface{}, messageQueueSize),
		processFunc:    processFunc,
	}
}

func (pool *Pool) PushMessage(message interface{}) {
	pool.messageChannel <- message
}

func (pool *Pool) Run() {
	for w := 1; w <= pool.capacity; w++ {
		pool.wg.Add(1)
		go pool.handleMessage(w, pool.messageChannel)
	}
}

func (pool *Pool) handleMessage(id int, channel <-chan interface{}) {
	defer pool.wg.Done()
	for message := range channel {
		fmt.Printf("worker %v is processing message %v\n", id, message)
		err := pool.processFunc(message)
		if err != nil {
			fmt.Printf("worker %v failed to process message %v with error: %v\n", id, message, err)
		}
	}
}

// Shutdown implement graceful shutdown
func (pool *Pool) Shutdown() {
	fmt.Printf("Closing worker pool... \n")
	close(pool.messageChannel)
	pool.wg.Wait()
	fmt.Printf("Closed worker pool!\n")
}
