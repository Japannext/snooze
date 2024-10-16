package utils

import (
	"sync"
)

type Semaphore struct {
	mu sync.Mutex
	cur int64
	max int64
}

func NewSemaphore(max int64) *Semaphore {
	return &Semaphore{
		cur: 0,
		max: max,
	}
}

func (sem *Semaphore) TryAcquire() bool {
	sem.mu.Lock()

	if sem.cur >= sem.max {
		return false
	}

	sem.cur += 1
	sem.mu.Unlock()
	return true
}

func (sem *Semaphore) Release() {
	sem.mu.Lock()
	if sem.cur > 1 {
		sem.cur -= 1
	}
	sem.mu.Unlock()
}

func (sem *Semaphore) Max() int64 {
	return sem.max
}

func (sem *Semaphore) Capacity() int64 {
	return sem.max - sem.cur
}

func (sem *Semaphore) WaitForFree() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		for {
			
		}
	}()
	return ch
}
