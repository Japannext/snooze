package service

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"

	"github.com/japannext/snooze/components/processor/decision"
	"github.com/japannext/snooze/internal/kv"
	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/pipeline"
)

type Processor struct {
	// A key-value that can quickly find which service a log
	// belong to based on labels.
	kv *kv.Kv[*models.Log, pipeline.Service]
}

func New(cfgs []pipeline.Service) (*Processor, error) {
	kv, err := kv.New[*models.Log](cfgs)
	if err != nil {
		return &Processor{}, fmt.Errorf("failed to initialize kv: %w", err)
	}

	return &Processor{
		kv: kv,
	}, nil
}

func (p *Processor) Process(ctx context.Context, item *models.Log) *decision.Decision {
	_, span := otel.Tracer("snooze").Start(ctx, "service")
	defer span.End()

	svc, found := p.kv.FirstMatch(item)
	if !found {
		return decision.OK()
	}

	item.Service = svc.Name

	return decision.OK()
}
