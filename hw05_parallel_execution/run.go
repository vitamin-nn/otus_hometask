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
func Run(tasks []Task, n, m int) error {
	if m < minM || n < minN {
		return ErrInvalidParams
	}
	if len(tasks) == 0 {
		return nil
	}
	done := make(chan struct{})

	workerCnt := n
	if len(tasks) < n {
		workerCnt = len(tasks)
	}
	taskCh := startProducer(tasks, done)
	errSignalCh := startWorkers(workerCnt, taskCh, done)
	reachErrLimit := startErrorCounter(m, errSignalCh, done)

	var err error
	if reachErrLimit {
		err = ErrErrorsLimitExceeded
	}
	return err
}

func startProducer(tasks []Task, done <-chan struct{}) <-chan Task {
	taskCh := make(chan Task)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, t := range tasks {
			select {
			case <-done:
				return
			default:
				taskCh <- t
			}
		}
	}()
	go func() {
		wg.Wait()
		close(taskCh)
	}()
	return taskCh
}

func startWorkers(n int, taskCh <-chan Task, done <-chan struct{}) <-chan struct{} {
	errSignalCh := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for t := range taskCh {
				select {
				case <-done:
					return
				default:
				}
				err := t()
				if err != nil {
					errSignalCh <- struct{}{}
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(errSignalCh)
	}()
	return errSignalCh
}

func startErrorCounter(m int, errorCh <-chan struct{}, done chan struct{}) bool {
	i := 0
	reachErrLimit := false
	for range errorCh {
		i++
		if i == m {
			close(done)
			reachErrLimit = true
		}
	}
	return reachErrLimit
}
