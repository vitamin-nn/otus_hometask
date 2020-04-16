package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

const (
	minM = 1
	minN = 1
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrInvalidParams = errors.New("invalid input params")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, n int, m int) error {
	if m < minM || n < minN {
		return ErrInvalidParams
	}
	if len(tasks) == 0 {
		return nil
	}

	taskCh := make(chan Task)
	errorCh := make(chan struct{})
	done := make(chan struct{})

	workerCnt := n
	if len(tasks) < n {
		workerCnt = len(tasks)
	}
	wgWorker := sync.WaitGroup{}
	wgWorker.Add(workerCnt)
	workerFunc := func() {
		defer wgWorker.Done()
		for t := range taskCh {
			select {
			case <-done:
				return
			default:
			}

			err := t()
			if err != nil {
				errorCh <- struct{}{}
			}
		}
	}

	producerFunc := func() {
		for _, t := range tasks {
			select {
			case <-done:
				return
			default:
				taskCh <- t
			}
		}
	}

	// errors counter
	reachErrLimit := false
	wgErrCounter := sync.WaitGroup{}
	wgErrCounter.Add(1)
	errCounterFunc := func() {
		defer wgErrCounter.Done()
		i := 0
		for range errorCh {
			i++
			if i == m {
				close(done)
				reachErrLimit = true
			}
		}
	}

	for i := 0; i < workerCnt; i++ {
		go workerFunc()
	}
	go errCounterFunc()
	producerFunc()

	close(taskCh)
	wgWorker.Wait()
	close(errorCh)
	wgErrCounter.Wait()

	var err error
	if reachErrLimit {
		err = ErrErrorsLimitExceeded
	}
	return err
}
