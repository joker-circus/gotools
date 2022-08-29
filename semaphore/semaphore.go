package semaphore

import ()

type Semaphore struct {
	bufSize int
	channel chan int8
}

func NewSemaphore(concurrencyNum int) *Semaphore {
	return &Semaphore{channel: make(chan int8, concurrencyNum), bufSize: concurrencyNum}
}

func (s *Semaphore) TryAcquire() bool {
	select {
	case s.channel <- int8(0):
		return true
	default:
		return false
	}
}

func (s *Semaphore) Acquire() {
	s.channel <- int8(0)
}

func (s *Semaphore) Release() {
	<-s.channel
}

func (s *Semaphore) AvailablePermits() int {
	return s.bufSize - len(s.channel)
}
