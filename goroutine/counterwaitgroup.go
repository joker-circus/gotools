package goroutine

import (
	"sync"
	"sync/atomic"
)

// CounterWaitGroup
type CounterWG struct {
	limit          int32
	currentRoutine int32
	wg             sync.WaitGroup
}

const defaultLimit = 10000

func CounterWaitGroup() *CounterWG {
	return &CounterWG{
		limit: defaultLimit,
	}
}

func (w *CounterWG) CurrentConcurrent() int32 {
	return atomic.LoadInt32(&w.currentRoutine)
}

func (w *CounterWG) LimitConcurrent() int32 {
	return atomic.LoadInt32(&w.limit)
}

func (w *CounterWG) SetLimit(limit int32) {
	if limit <= 0 {
		return
	}
	atomic.StoreInt32(&w.limit, limit)
}

// 如果影响超过上限，则会返回 false
func (w *CounterWG) Go(cb func()) bool {

	if w.LimitConcurrent() <= w.CurrentConcurrent() {
		return false
	}

	atomic.AddInt32(&w.currentRoutine, 1)
	w.wg.Add(1)

	go func() {
		cb()
		w.wg.Done()
		atomic.AddInt32(&w.currentRoutine, -1)
	}()
	return true
}

func (w *CounterWG) Wait() {
	w.wg.Wait()
}
