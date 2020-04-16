package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, (func(i int) Task {
				return func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					atomic.AddInt32(&runTasksCount, 1)
					return fmt.Errorf("error from task %d", i)
				}
			})(i))
		}

		workersCount := 10
		maxErrorsCount := 23
		result := Run(tasks, workersCount, maxErrorsCount)
		require.Equal(t, ErrErrorsLimitExceeded, result)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		result := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.Nil(t, result)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("tasks count are less than gorutines", func(t *testing.T) {
		tasksCount := 6
		taskErrCount := 2
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		var task Task
		// adding without error tasks
		for i := 0; i < tasksCount; i++ {
			task = func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			}
			tasks = append(tasks, task)
		}

		// adding error tasks
		for i := 0; i < taskErrCount; i++ {
			task = func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return fmt.Errorf("error from task %d", i)
			}
			tasks = append(tasks, task)
		}

		workersCount := 10
		maxErrorsCount := 2

		result := Run(tasks, workersCount, maxErrorsCount)
		require.Equal(t, ErrErrorsLimitExceeded, result)

		require.Equal(t, int32(tasksCount+taskErrCount), runTasksCount, "not all tasks were completed")
	})

	t.Run("border conditions check", func(t *testing.T) {
		tasks := make([]Task, 0, 1)

		result := Run(tasks, 10, 1)
		require.Nil(t, result)

		task := func() error {
			return nil
		}
		tasks = append(tasks, task)
		result = Run(tasks, 10, 0)
		require.Equal(t, ErrInvalidParams, result)

		result = Run(tasks, 0, 1)
		require.Equal(t, ErrInvalidParams, result)
	})
}
