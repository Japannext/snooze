package utils

import (
	"time"
)

type Batch[T any] struct {
	Items []T `json:"items"`
}

type BatchingChannel[T any] struct {
	in chan T
	out chan Batch[T]
	size int
	timeout time.Duration
	stopping bool
	done (chan bool)
}

func NewBatchingChannel[T any](size int, timeout time.Duration) *BatchingChannel[T] {
	bc := &BatchingChannel[T]{
		in: make(chan T),
		out: make(chan Batch[T]),
		size: size,
		timeout: timeout,
	}
	bc.Start()
	return bc
}

func newBatch[T any]() Batch[T] {
	return Batch[T]{Items: []T{}}
}

func (bc *BatchingChannel[T]) Start() {
	go func() {
		var batch = newBatch[T]()
		var ticker = time.NewTicker(bc.timeout)
		defer ticker.Stop()
		for {
			if bc.stopping {
				break
			}
			select {
			case item := <-bc.in:
				batch.Items = append(batch.Items, item)
				if len(batch.Items) >= bc.size {
					bc.out <- batch
					batch = newBatch[T]()
					ticker.Reset(bc.timeout)
					continue
				}
			case <-ticker.C: // timer
				if len(batch.Items) > 0 {
					bc.out <- batch
					batch = newBatch[T]()
					continue
				}
			}
		}
		bc.done <-true
	}()
}

func (bc *BatchingChannel[T]) Stop() {
	bc.stopping = true
	<-bc.done
	close(bc.in)
	close(bc.out)
}

func (bc *BatchingChannel[T]) Add(item T) {
	bc.in <- item
}

func (bc *BatchingChannel[T]) Channel() (<-chan Batch[T]) {
	return bc.out
}
