package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

var TempErr int32
var NumberTasks int32

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	var idx int32
	var tempErr = int32(m)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				temp := atomic.AddInt32(&idx, 1)
				if tempErr >= 0 && int(temp-1) < len(tasks) {
					err := tasks[temp-1]()
					atomic.AddInt32(&NumberTasks, 1)
					if err != nil {
						atomic.AddInt32(&tempErr, -1)
					}
				} else {
					return
				}

			}
		}()
	}
	wg.Wait()
	if tempErr < 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}
