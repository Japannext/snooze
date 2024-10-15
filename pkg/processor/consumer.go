package processor

import (
	"context"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/japannext/snooze/pkg/models"

	"github.com/japannext/snooze/pkg/processor/activecheck"
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
		batch, err := processQ.Fetch(config.BatchSize, jetstream.FetchMaxWait(5*time.Second))
		if err != nil {
			log.Warnf("failed to fetch items: %s", err)
			continue
		}
		var msgs []jetstream.Msg
		for msg := range batch.Messages() {
			msgs = append(msgs, msg)
		}
		log.Debugf("Fetched %d logs", len(msgs))
		if len(msgs) == 0 {
			continue
		}
		ctx := context.TODO()
		processBatch(ctx, msgs)
	}
}

func (c *Consumer) Stop() {
}

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

func processBatch(ctx context.Context, msgs []jetstream.Msg) {

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

	processedLogs.Inc()
}
