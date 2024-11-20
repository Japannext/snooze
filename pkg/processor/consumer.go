package processor

import (
	"context"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/tracing"

	"github.com/japannext/snooze/pkg/processor/activecheck"
	"github.com/japannext/snooze/pkg/processor/grouping"
	"github.com/japannext/snooze/pkg/processor/notification"
	"github.com/japannext/snooze/pkg/processor/profile"
	"github.com/japannext/snooze/pkg/processor/ratelimit"
	"github.com/japannext/snooze/pkg/processor/silence"
	"github.com/japannext/snooze/pkg/processor/snooze"
	"github.com/japannext/snooze/pkg/processor/store"
	"github.com/japannext/snooze/pkg/processor/timestamp"
	"github.com/japannext/snooze/pkg/processor/transform"
	"github.com/japannext/snooze/pkg/processor/mapping"
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
		msgs, err := processQ.Fetch(size)
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
		item = models.Log{
			ActualTime: models.TimeNow(),
			DisplayTime: models.TimeNow(),
			Identity: map[string]string{"snooze.internal": "error"},
			Message: string(data),
			Error: err.Error(),
		}
	}
	return &item
}

/*
type ProcessDecision int
const (
	CONTINUE ProcessDecision = iota
	SILENCE
	DROP
)

type ProcessFunc = func(context.Context, *models.Log)

var processes = []ProcessFunc{
	timestamp.Process,
	transform.Process,
	silence.Process,
	profile.Process,
	grouping.Process,
	ratelimit.Process,
	activecheck.Process,
	snooze.Process,
	notification.Process,
	store.Process,
}
*/

func processLog(ctx context.Context, item *models.Log) {
	ctx, span := tracer.Start(ctx, "processLog")
	defer span.End()

	start := time.Now()

	tracing.SetLog(span, item)

	ctx = mapping.WithMappings(ctx)

	// Transformative
	timestamp.Process(ctx, item)
	transform.Process(ctx, item)
	profile.Process(ctx, item)
	grouping.Process(ctx, item)

	// Traffic control
	silence.Process(ctx, item)
	ratelimit.Process(ctx, item)

	// Snooze
	snooze.Process(ctx, item)

	// Monitoring
	activecheck.Process(ctx, item)

	// Notification + Storage
	notification.Process(ctx, item)
	store.Process(ctx, item)

	processedLogs.Inc()
	processTime.Observe(time.Since(start).Seconds())
}
