package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/rabbitmq"

	"github.com/japannext/snooze/pkg/processor/activecheck"
	"github.com/japannext/snooze/pkg/processor/grouping"
	"github.com/japannext/snooze/pkg/processor/profile"
	"github.com/japannext/snooze/pkg/processor/notification"
	"github.com/japannext/snooze/pkg/processor/ratelimit"
	"github.com/japannext/snooze/pkg/processor/silence"
	"github.com/japannext/snooze/pkg/processor/store"
	"github.com/japannext/snooze/pkg/processor/transform"
)

type Processor struct{
	Consumer *rabbitmq.Consumer
}

func NewProcessor() *Processor {
	processor := &Processor{}

	rabbitmq.SetupProcessing()

	options := rabbitmq.ConsumerOptions{}
	processor.Consumer = rabbitmq.NewConsumer(rabbitmq.PROCESSING_QUEUE, "", handler, options)
	return processor
}

// For item that will not be requeued, because their
// format is invalid, or they are poison messages.
type RejectedLog struct {
	item  *api.Log
	reason string
}

func (r *RejectedLog) Error() string {
	// return fmt.Sprintf("Rejected item id=%s/%s reason=%s", r.item.TraceID, r.item.SpanID, r.reason)
	return fmt.Sprintf("Rejected item: %s", r.reason)
}

func (p *Processor) Run() error {
	if err := p.Consumer.ConsumeForever(); err != nil {
		return err
	}
	return nil
}

func (p *Processor) Stop() {
	p.Consumer.Stop()
}

func handler(delivery rabbitmq.Delivery) error {
	var item *api.Log
	if err := json.Unmarshal(delivery.Body, &item); err != nil {
		return err
	}
	if err := Process(item); err != nil {
		return err
	}
	delivery.Ack(false)
	return nil
}

func Process(item *api.Log) error {
	ctx := context.TODO()
	start := time.Now()

	log.Debugf("Start processing item: %s", item)
	if err := transform.Process(ctx, item); err != nil {
		return err
	}
	if err := silence.Process(ctx, item); err != nil {
		return err
	}
	if err := profile.Process(ctx, item); err != nil {
		return err
	}
	if err := grouping.Process(ctx, item); err != nil {
		return err
	}
	if err := ratelimit.Process(ctx, item); err != nil {
		return err
	}

	// Active-check
	activecheck.Process(ctx, item)

	if err := notification.Process(ctx, item); err != nil {
		return err
	}
	if err := store.Process(ctx, item); err != nil {
		return err
	}
	log.Debugf("End processing item: %s", item)

	processMetrics(start, item)

	return nil
}
