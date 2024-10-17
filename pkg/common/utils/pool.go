package utils

import (
	"sync"
)

type Pool struct {
	busy chan struct{}
	max int
	ready *sync.Cond
}

func NewPool(max int) *Pool {
	mu := sync.Mutex{}
	pool := &Pool{
		max: max,
		busy: make(chan struct{}, max),
		ready: sync.NewCond(&mu),
	}
	return pool
}

func (pool *Pool) Ready() int {
	pool.ready.L.Lock()
	defer pool.ready.L.Unlock()
	if len(pool.busy) == pool.max {
		pool.ready.Wait()
	}
	return pool.max - len(pool.busy)
}

func (pool *Pool) TryAcquire() bool {
	pool.ready.L.Lock()
	defer pool.ready.L.Unlock()
	select {
	case pool.busy <-struct{}{}:
		return true
	default:
		return false
	}

}

func (pool *Pool) Max() int {
	return pool.max
}

func (pool *Pool) Busy() int {
	return len(pool.busy)
}

func (pool *Pool) Release() {
	pool.ready.L.Lock()
	defer pool.ready.L.Unlock()
	select { // Release without waiting
	case <-pool.busy:
		pool.ready.Signal()
	default:
	}
}
