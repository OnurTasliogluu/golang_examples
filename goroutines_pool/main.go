package main

import (
	"fmt"
	"sync"
	"time"

	"goroutines_pool/m/v2/pkg/workpool"
)

func main() {
	// Create a new workpool with 3 max concurrent workers
	pool := workpool.New(3)

	// Shared resource to demonstrate concurrent work
	var counter int
	var mu sync.Mutex

	// Insert tasks
	for i := 0; i < 10; i++ {
		taskNum := i
		pool.Insert(func() {
			// Simulate some work
			time.Sleep(time.Second)

			mu.Lock()
			counter++
			fmt.Printf("Task %d completed. Counter: %d\n", taskNum, counter)
			mu.Unlock()
		})
	}

	// Run all tasks and wait for completion
	pool.RunAndWait()

	fmt.Printf("All tasks finished. Final counter: %d\n", counter)
}
