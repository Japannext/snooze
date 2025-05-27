package store

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"

	"github.com/japannext/snooze/components/processor/decision"
	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/mq"
	"github.com/japannext/snooze/internal/opensearch/format"
	"github.com/japannext/snooze/internal/tracing"
)

type Processor struct {
	storeQ *mq.Pub
}

func New() (*Processor, error) {
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
