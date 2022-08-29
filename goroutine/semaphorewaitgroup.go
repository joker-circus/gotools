package goroutine

import (
	"context"
	"sync"

	mysema "github.com/joker-circus/gotools/semaphore"
	gosema "golang.org/x/sync/semaphore"
)

type Semaphore interface {
	Acquire()
	Release()
}

// 简易版信号量控制的 WaitGroup
func SimpleSemaphoreWaitGroup(limit int) *SemaphoreWG {
	return SemaphoreWaitGroup(mysema.NewSemaphore(limit))
}

// Go 官方权重信号量控制的 WaitGroup
func GoSemaphoreWaitGroup(limit int64) *SemaphoreWG {
	return SemaphoreWaitGroup(&goSemaphore{gosema.NewWeighted(limit)})
}

// Semaphore WaitGroup
type SemaphoreWG struct {
	sema Semaphore
	wg   sync.WaitGroup
}

func SemaphoreWaitGroup(semaphore Semaphore) *SemaphoreWG {
	return &SemaphoreWG{
		sema: semaphore,
	}
}

func (w *SemaphoreWG) Go(cb func()) {
	w.sema.Acquire()
	w.wg.Add(1)

	go func() {
		cb()
		w.wg.Done()
		w.sema.Release()
	}()
}

func (w *SemaphoreWG) Wait() {
	w.wg.Wait()
}

type goSemaphore struct {
	w *gosema.Weighted
}

func (gs *goSemaphore) Acquire() {
	_ = gs.w.Acquire(context.Background(), 1)
}

func (gs *goSemaphore) Release() {
	gs.w.Release(1)
}
