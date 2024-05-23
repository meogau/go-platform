package main

import (
	"context"
	"fmt"
	"github.com/meogau/go-platform/worker"
	"os/signal"
	"syscall"
)

func main() {
	//Test worker pool
	workerPool := worker.NewWorkerPool(5, 100, func(message interface{}) error {
		fmt.Printf("processed message %v", message)
		return nil
	})
	workerPool.Run()
	for i := 0; i < 20; i++ {
		workerPool.PushMessage(fmt.Sprintf("message %v", i))
	}
	//graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	fmt.Println("Receive close signal")
	workerPool.Shutdown()
}
