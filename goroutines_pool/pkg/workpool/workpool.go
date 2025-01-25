package workpool

import (
	"container/list"
	"sync"
)

// Task represents a function to be executed
type Task func()

// WorkPool manages a pool of tasks to be executed
type WorkPool struct {
	tasks      *list.List
	mu         sync.Mutex
	maxWorkers int
	wg         sync.WaitGroup
}

// New creates a new WorkPool with specified max workers
func New(maxWorkers int) *WorkPool {
	return &WorkPool{
		tasks:      list.New(),
		maxWorkers: maxWorkers,
	}
}

// Insert adds a task to the end of the queue
func (wp *WorkPool) Insert(task Task) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.tasks.PushBack(task)
}

// InsertPosition adds a task at a specific position in the queue
func (wp *WorkPool) InsertPosition(task Task, position int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if position < 0 {
		wp.tasks.PushFront(task)
		return
	}

	current := wp.tasks.Front()
	for i := 0; i < position && current != nil; i++ {
		current = current.Next()
	}

	if current != nil {
		wp.tasks.InsertBefore(task, current)
	} else {
		wp.tasks.PushBack(task)
	}
}

// List returns a copy of all tasks in the queue
func (wp *WorkPool) List() []Task {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	tasks := make([]Task, 0, wp.tasks.Len())
	for e := wp.tasks.Front(); e != nil; e = e.Next() {
		tasks = append(tasks, e.Value.(Task))
	}
	return tasks
}

// Remove removes a task at the specified index
func (wp *WorkPool) Remove(index int) bool {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if index < 0 {
		return false
	}

	current := wp.tasks.Front()
	for i := 0; i < index && current != nil; i++ {
		current = current.Next()
	}

	if current != nil {
		wp.tasks.Remove(current)
		return true
	}
	return false
}

// RunAndWait executes all tasks concurrently and waits for completion
func (wp *WorkPool) RunAndWait() {
	wp.mu.Lock()
	taskCount := wp.tasks.Len()
	wp.mu.Unlock()

	if taskCount == 0 {
		return
	}

	// Create a semaphore channel to limit concurrent workers
	sem := make(chan struct{}, wp.maxWorkers)

	for wp.tasks.Len() > 0 {
		wp.mu.Lock()
		taskElement := wp.tasks.Front()
		wp.tasks.Remove(taskElement)
		task := taskElement.Value.(Task)
		wp.mu.Unlock()

		wp.wg.Add(1)
		sem <- struct{}{}

		go func(t Task) {
			defer wp.wg.Done()
			defer func() { <-sem }()
			t()
		}(task)
	}

	wp.wg.Wait()
}
