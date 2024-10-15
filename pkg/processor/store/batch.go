package store

import (
	"context"
	"sync"
	"time"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/processor/tracing"
)

var QUEUE = fmt.Sprintf("%s.%s", mq.STORE_STREAM, models.LOG_INDEX)

func Batch(ctx context.Context, msgs []jetstream.Msg, items []*models.Log) error {
	ctx, span := tracing.TRACER.Start(ctx, "store")
	defer span.End()

	var wg sync.WaitGroup
	for i, item := range items {
		if item.Mute.SkipStorage {
			continue
		}
		wg.Add(1)
		msg := msgs[i]
		future, err := mq.PublishAsync(QUEUE, item)
		if err != nil {
			log.Warnf("failed to publish log: %s", err)
			msg.Nak()
			continue
		}
		wg.Add(1)
		go func() {
			select {
			case <-future.Ok():
				msg.Ack()
			case <-future.Err():
				log.Warnf("failed to publish log (async): %s", err)
				msg.Ack()
			case <-time.After(5*time.Second):
				log.Warnf("failed to publish log (async): timeout")
				msg.Nak()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	return nil
}
