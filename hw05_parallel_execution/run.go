package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrInvalidArgument = errors.New("invalid argument")

type Task func() error

type workerResults struct {
	error error
	task  Task
}

func Run(tasks []Task, poolSize int, maxErrors int) error {
	var exitError error

	if poolSize < 0 || maxErrors < 0 {
		return ErrInvalidArgument
	}

	if len(tasks) == 0 {
		return nil
	}

	var queue = make(chan Task, len(tasks))
	// push to queue
	for _, task := range tasks {
		queue <- task
	}

	var wgWorkers sync.WaitGroup
	var doneWorkers = make(chan bool)
	var workersResults = make(chan *workerResults, len(tasks))
	// run workers
	for i := 0; i < poolSize; i++ {
		wgWorkers.Add(1)
		go func() {
			defer wgWorkers.Done()
			worker(queue, workersResults, doneWorkers)
		}()
	}

	// handle workers results
	var qtyErrors int
	var qtyNoErrorTasks int
	var attempts int
	var maxAttempts = poolSize + maxErrors
	for result := range workersResults {
		attempts++
		// complete all tasks
		if result.error != nil {
			qtyErrors++
			// return to queue if error
			queue <- result.task
		} else {
			qtyNoErrorTasks++
		}

		if qtyNoErrorTasks >= len(tasks) {
			break
			// if errors then all attempts = poolSize + maxErrors
		} else if qtyErrors > 0 && attempts >= maxAttempts {
			exitError = ErrErrorsLimitExceeded
			break
		}
	}

	close(doneWorkers)
	wgWorkers.Wait()

	close(queue)
	close(workersResults)

	return exitError
}

func worker(queue <-chan Task, result chan<- *workerResults, done <-chan bool) {
	handler := func(task Task) {
		result <- &workerResults{error: task(), task: task}
	}

	for {
		select {
		case task, ok := <-queue:
			if ok {
				handler(task)
			}
		case <-done:
			return
		}
	}
}
