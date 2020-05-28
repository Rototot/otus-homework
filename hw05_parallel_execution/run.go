package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type workerResults struct {
	error error
	task  Task
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, N int, M int) error {
	var queue = make(chan Task, len(tasks))
	// push to queue
	for _, task := range tasks {
		queue <- task
	}

	var wgWorkers sync.WaitGroup
	var doneWorkers = make(chan bool)
	var workersResults = make(chan *workerResults, N)
	// run workers
	for i := 0; i < N; i++ {
		wgWorkers.Add(1)
		go func() {
			defer wgWorkers.Done()
			worker(queue, workersResults, doneWorkers)
		}()
	}

	// results handler
	var qtyErrors int
	var qtyNoErrorTasks int
	var attempts int
	var maxAttempts = len(tasks) + M
	for result := range workersResults {
		// complete all tasks
		// or
		// if errors then all attempts = N + M
		if qtyNoErrorTasks >= len(tasks) || qtyErrors > 0 && attempts >= maxAttempts {
			close(doneWorkers)
			break
		}

		attempts++
		if result.error != nil {
			qtyErrors++
			// return to queue
			queue <- result.task
		} else {
			qtyNoErrorTasks++
		}
	}

	wgWorkers.Wait()

	close(queue)
	close(workersResults)

	if qtyErrors > 0 && attempts >= maxAttempts {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(queue <-chan Task, result chan<- *workerResults, done <-chan bool) {

	handler := func(task Task) {
		result <- &workerResults{error: task(), task: task}
	}

	for {
		select {
		case task := <-queue:
			handler(task)
		case <-done:
			return
		}
	}
}
