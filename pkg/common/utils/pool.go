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
	pool := &Pool{
		max: max,
		busy: make(chan struct{}, max),
		ready: sync.NewCond(new(sync.Mutex)),
	}
	return pool
}

func (pool *Pool) Ready() int {
	if len(pool.busy) == pool.max {
		pool.ready.Wait()
	}
	return pool.max - len(pool.busy)
}

func (pool *Pool) TryAcquire() bool {
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
	select { // Release without waiting
	case <-pool.busy:
	default:
	}
	pool.ready.Signal()
}
