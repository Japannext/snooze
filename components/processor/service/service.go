package service

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"

	"github.com/japannext/snooze/components/processor/decision"
	"github.com/japannext/snooze/internal/kv"
	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/pipeline"
	"github.com/japannext/snooze/internal/tracing"
)

type Processor struct {
	// A key-value that can quickly find which service a log
	// belong to based on labels.
	kv *kv.Kv[pipeline.Service]
}

func New(cfgs []pipeline.Service) (*Processor, error) {
	kv := kv.New[pipeline.Service]()
	for _, cfg := range cfgs {
		for _, selector := range cfg.Selectors {
			if err := kv.Add(selector, &cfg); err != nil {
				return &Processor{}, fmt.Errorf("failed in service %s: %w", cfg.Name, err)
			}
		}
	}

	return &Processor{
		kv: kv,
	}, nil
}

func (p *Processor) Process(ctx context.Context, item *models.Log) *decision.Decision {
	_, span := otel.Tracer("snooze").Start(ctx, "service")
	defer span.End()

	for _, h := range item.SelectorHashes {
		svc, found := p.kv.FirstMatch(h)
		if found {
			item.Service = svc.Name
			tracing.SetString(span, "serviceName", svc.Name)
			return decision.OK()
		}
	}

	// Not found
	return decision.OK()
}
