package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrInvalidArgument = errors.New("invalid argument")

type Task func() error

func Run(tasks []Task, poolSize int, maxErrors int) error {
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

	// run pool workers
	var pool = NewWorkerPool(queue, poolSize, maxErrors, poolSize+maxErrors)
	err := pool.Listen()

	close(queue)

	return err
}

type workersPool struct {
	queue           chan Task
	size            int
	maxErrors       int
	maxAttempts     int
	maxNoErrorTasks int
	qtyErrors       int
	qtyAttempts     int
	qtyNoErrorTasks int
	su              sync.Mutex
}

func NewWorkerPool(queue chan Task, poolSize int, maxErrors int, maxAttempts int) *workersPool {
	return &workersPool{
		queue:           queue,
		size:            poolSize,
		maxErrors:       maxErrors,
		maxAttempts:     maxAttempts,
		maxNoErrorTasks: cap(queue),
	}
}

func (p *workersPool) Listen() error {
	// await complete workers
	p.startWorkers()

	if p.qtyErrors > p.maxErrors {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func (p *workersPool) startWorkers() {
	var wg sync.WaitGroup
	for i := 0; i < p.size; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.worker()
		}()
	}

	wg.Wait()
}

func (p *workersPool) worker() {
	var isNeedStop = func() bool {
		p.su.Lock()
		defer p.su.Unlock()

		// errors > M
		// with erros -> >= N+M
		// complete = N
		return p.qtyErrors > p.maxErrors ||
			(p.qtyErrors > 0 && p.qtyAttempts >= p.maxAttempts) ||
			p.qtyNoErrorTasks >= p.maxNoErrorTasks
	}
	for {
		select {
		case task, ok := <-p.queue:
			if !ok || isNeedStop() {
				return
			}
			p.handleTask(task)
		default:
			if isNeedStop() {
				return
			}
		}
	}
}

func (p *workersPool) handleTask(task Task) {
	err := task()
	p.su.Lock()
	defer p.su.Unlock()

	p.qtyAttempts++
	if err == nil {
		p.qtyNoErrorTasks++
	} else {
		p.qtyErrors++
	}
}