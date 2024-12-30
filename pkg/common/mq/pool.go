package mq

import (
	"context"

	"github.com/japannext/snooze/pkg/common/utils"
	"github.com/nats-io/nats.go/jetstream"
)

/* A daemon requiring a handler that auto manage a pool of workers
while consuming a queue. */

type MessageHandler = func(context.Context, jetstream.Msg) error

type WorkerPool struct {
	q        *Sub
	handler  MessageHandler
	pool     *utils.Pool
	stopping bool
	done     chan struct{}
}

func NewWorkerPool(sub *Sub, handler MessageHandler, workers int) *WorkerPool {
	wp := &WorkerPool{
		q:        sub,
		handler:  handler,
		pool:     utils.NewPool(workers),
		stopping: false,
		done:     make(chan struct{}),
	}

	return wp
}

var RETRY = &Retry{}

type Retry struct{}

func (r *Retry) Error() string {
	return "retry"
}

func (wp *WorkerPool) Run() error {
	for {
		if wp.stopping {
			break
		}
		size := wp.pool.Ready()
		msgs, err := wp.q.Fetch(size)
		if err != nil {
			log.Warnf("failed to fetch items: %s", err)
			continue
		}
		if len(msgs) == 0 {
			continue
		}
		for _, m := range msgs {
			if ok := wp.pool.TryAcquire(); !ok {
				log.Warnf("no worker ready, requeuing!")
				m.Msg.Nak()
				continue
			}
			go func() {
				defer wp.pool.Release()
				err := wp.handler(m.Context, m.Msg)
				if err == RETRY {
					m.Msg.Nak()
					return
				}
				if err != nil {
					m.Msg.Term()
					return
				}
				m.Msg.Ack()
			}()
		}
	}
	wp.done <- struct{}{}
	return nil
}

func (wp *WorkerPool) Stop() {
	wp.stopping = true
	<-wp.done
}
