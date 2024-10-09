package processor

import (
	"context"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/rabbitmq"

	"github.com/japannext/snooze/pkg/processor/activecheck"
	"github.com/japannext/snooze/pkg/processor/grouping"
	"github.com/japannext/snooze/pkg/processor/profile"
	"github.com/japannext/snooze/pkg/processor/notification"
	"github.com/japannext/snooze/pkg/processor/ratelimit"
	"github.com/japannext/snooze/pkg/processor/silence"
	"github.com/japannext/snooze/pkg/processor/store"
	"github.com/japannext/snooze/pkg/processor/timestamp"
	"github.com/japannext/snooze/pkg/processor/transform"
	"github.com/japannext/snooze/pkg/processor/tracing"
	"github.com/japannext/snooze/pkg/common/utils"
)

type Process interface {
	Name() string
	Batch(models.BatchOfLogs) error
}

type Processor struct {
	processes []BatchFunc
	consumer *rabbitmq.Consumer
}

func (processor *Processor) Batch(batch models.BatchOfLogs) error {
	for _, p := range processor.processes {
		if err := p(batch); err != nil {
			log.Warnf("[process:%s] %s", p.Name(), err)
			continue
		}
	}
}

type BatchFunc = func(context.Context, []*models.Log) error

func NewProcessor() *Processor {
	processor := &Processor{
		processes: []Process{
			transform.Batch,
			silence.Batch,
			profile.Batch,
			grouping.Batch,
			ratelimit.Batch,
			activecheck.Batch,
			notification.Batch,
			store.Batch,
		},
	}

	rabbitmq.SetupProcessing()

	options := rabbitmq.ConsumerOptions{}
	processor.Consumer = rabbitmq.NewBatchConsumer(rabbitmq.PROCESSING_QUEUE, "", batchHandler, options)
	return processor
}

func (p *Processor) Run() error {
	//if err := p.Consumer.ConsumeForever(); err != nil {
	timeout := time.Duration(config.BatchTimeoutSeconds) * time.Second
	if err := p.Consumer.BatchConsume(config.BatchSize, timeout); err != nil {
		return err
	}
	return nil
}

func (p *Processor) Stop() {
	p.Consumer.Stop()
}

func batchHandler(deliveries []rabbitmq.Delivery) error {
	for _, delivery := range deliveries {
		var batch utils.Batch[*models.Log]
		if err := json.Unmarshal(delivery.Body, &batch); err != nil {
			log.Warnf("cannot unmarshal `%s`: %s", delivery.Body, err)
			delivery.Reject(false)
			continue
		}
		if err := Batch(batch); err != nil {
			delivery.Reject(false)
			return err
		}
		delivery.Ack(false)
	}
	return nil
}

func Batch(batch utils.Batch[*models.Log]) error {
	ctx := context.TODO()
	start := time.Now()

	// In-queue time statistics
	ts := time.UnixMilli(int64(batch.TimestampMillis))
	inqueueTime.Observe(time.Since(ts).Seconds())

	// Tracing
	ctx, span := tracing.TRACER.Start(ctx, "process")
	defer span.End()

	items := batch.Items

	timestamp.Batch(ctx, items)

	log.Debugf("Start processing batch of %d items", len(items))
	if err := transform.Batch(ctx, items); err != nil {
		return err
	}
	if err := silence.Batch(ctx, items); err != nil {
		return err
	}
	if err := profile.Batch(ctx, items); err != nil {
		return err
	}
	if err := grouping.Batch(ctx, items); err != nil {
		return err
	}
	if err := ratelimit.Batch(ctx, items); err != nil {
		return err
	}

	// Active-check
	activecheck.Batch(ctx, items)

	if err := notification.Batch(ctx, items); err != nil {
		return err
	}
	if err := store.Batch(ctx, items); err != nil {
		return err
	}
	log.Debugf("End processing %d items", len(items))

	batchTime.Observe(time.Since(start).Seconds())
	batchSize.Observe(float64(len(items)))
	processedLogs.Add(float64(len(batch.Items)))

	return nil
}
