package store

import (
	"context"

	"github.com/japannext/snooze/components/processor/decision"
	"github.com/japannext/snooze/internal/mq"
	"github.com/japannext/snooze/internal/opensearch/format"
	"github.com/japannext/snooze/internal/tracing"
	"github.com/japannext/snooze/internal/models"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type Processor struct {
	storeQ *mq.Pub
}

type Config struct {
}

func New(cfg Config) (*Processor, error) {
	p := &Processor{}
	p.storeQ = mq.StorePub()

	return p, nil
}

func (p *Processor) Process(ctx context.Context, item *models.Log) *decision.Decision {
	ctx, span := otel.Tracer("snooze").Start(ctx, "store")
	defer span.End()

	tracing.SetBool(span, "status.skipStorage", item.Status.SkipStorage)

	if item.Status.SkipStorage {
		log.Debugf("skipping storage")

		return decision.OK()
	}

	err := p.storeQ.PublishData(ctx, &format.Create{
		Index: models.LogIndex,
		Item:  item,
	})
	if err != nil {
		return decision.Retry(err)
	}

	return decision.OK()
}
