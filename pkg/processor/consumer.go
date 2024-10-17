package processor

import (
	"context"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/japannext/snooze/pkg/models"

	"github.com/japannext/snooze/pkg/processor/grouping"
	"github.com/japannext/snooze/pkg/processor/profile"
	"github.com/japannext/snooze/pkg/processor/notification"
	"github.com/japannext/snooze/pkg/processor/ratelimit"
	"github.com/japannext/snooze/pkg/processor/silence"
	"github.com/japannext/snooze/pkg/processor/store"
	"github.com/japannext/snooze/pkg/processor/timestamp"
	"github.com/japannext/snooze/pkg/processor/transform"
)

type Consumer struct {}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Run() error {
	for {
		log.Debugf("Waiting for a worker to be ready...")
		size := pool.Ready()
		log.Debugf("%d workers are ready", size)
		msgs, err := processQ.Fetch(size, jetstream.FetchMaxWait(time.Duration(config.BatchTimeoutSeconds) * time.Second))
		if err != nil {
			log.Warnf("failed to fetch items: %s", err)
			continue
		}
		log.Debugf("Fetched %d logs", len(msgs))
		if len(msgs) == 0 {
			continue
		}

		for _, m := range msgs {
			msg := m.Msg
			// Immediately requeue messages when there is no worker
			if ok := pool.TryAcquire(); !ok {
				log.Warnf("Requeuing since no worker is ready (%d/%d)", pool.Busy(), pool.Max())
				msg.Nak()
				continue
			}
			go func() {
				defer pool.Release()
				processLog(m.Context, unmarshalLog(msg))
				msg.Ack()
			}()
		}
	}
}

func (c *Consumer) Stop() {
}

func unmarshalLog(msg jetstream.Msg) *models.Log {
	var item models.Log
	data := msg.Data()
	if err := json.Unmarshal(data, &item); err != nil {
		log.Warnf("invalid JSON while parsing log: %s", err)
		now := uint64(time.Now().UnixMilli())
		item = models.Log{
			Timestamp: models.Timestamp{Observed: now, Display: now},
			Identity: map[string]string{"snooze.internal": "error"},
			Message: string(data),
			Error: err.Error(),
		}
	}
	return &item
}

/*
func bytesToLogs(msgs []jetstream.Msg) []*models.Log {
	var items []*models.Log
	for _, msg := range msgs {
		var item models.Log
		data := msg.Data()
		if err := json.Unmarshal(data, &item); err != nil {
			log.Warnf("invalid JSON while parsing log: %s", err)
			now := uint64(time.Now().UnixMilli())
			item = models.Log{
				Timestamp: models.Timestamp{Observed: now, Display: now},
				Identity: map[string]string{"snooze.internal": "error"},
				Message: string(data),
				Error: err.Error(),
			}
		}
		items = append(items, &item)
	}
	return items
}
*/

func processLog(ctx context.Context, item *models.Log) {
	ctx, span := tracer.Start(ctx, "processLog")
	defer span.End()

	timestamp.Process(ctx, item)
	transform.Process(ctx, item)
	silence.Process(ctx, item)
	profile.Process(ctx, item)
	grouping.Process(ctx, item)
	ratelimit.Process(ctx, item)
	// activecheck.Process(ctx, item)
	notification.Process(ctx, item)
	store.Process(ctx, item)
}

/*
func processBatch(msgs []mq.MsgWithContext) {
	log.Debugf("processing %d logs", len(msgs))

	items := bytesToLogs(msgs)

	timestamp.Batch(ctx, items)
	transform.Batch(ctx, items)
	silence.Batch(ctx, items)
	profile.Batch(ctx, items)
	grouping.Batch(ctx, items)
	ratelimit.Batch(ctx, items)
	activecheck.Batch(ctx, items)
	notification.Batch(ctx, items)
	store.Batch(ctx, msgs, items)

	log.Debugf("waiting for publish queue to be exhausted...")
	<- mq.PublishAsyncComplete()
	log.Debugf("done publishing")

	for _, msg := range msgs {
		msg.Ack()
	}

	processedLogs.Add(float64(len(msgs)))
}
*/
