package processor

import (
	"context"
	"fmt"

	"github.com/japannext/snooze/components/processor/activecheck"
	"github.com/japannext/snooze/components/processor/decision"
	"github.com/japannext/snooze/components/processor/grouping"
	"github.com/japannext/snooze/components/processor/notification"
	"github.com/japannext/snooze/components/processor/problem"
	"github.com/japannext/snooze/components/processor/profile"
	"github.com/japannext/snooze/components/processor/ratelimit"
	"github.com/japannext/snooze/components/processor/service"
	"github.com/japannext/snooze/components/processor/silence"
	"github.com/japannext/snooze/components/processor/snooze"
	"github.com/japannext/snooze/components/processor/store"
	"github.com/japannext/snooze/components/processor/timestamp"
	"github.com/japannext/snooze/components/processor/transform"
	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/pipeline"
)

type SubProcessor interface {
	//Name() string
	Process(ctx context.Context, item *models.Log) *decision.Decision
}

// An object that can process a log, and contain the static config
// of all processors
type Processor struct {
	processes []SubProcessor
}

func NewProcessor(pl *pipeline.Pipeline) (*Processor, error) {
	p := &Processor{}

	transformProcess, err := transform.New(pl.Transforms)
	if err != nil {
		return p, fmt.Errorf("error in `transform`: %w", err)
	}

	profileProcess, err := profile.New(pl.Profiles)
	if err != nil {
		return p, fmt.Errorf("error in `profile`: %w", err)
	}

	groupingProcess, err := grouping.New(pl.Groupings)
	if err != nil {
		return p, fmt.Errorf("error in `grouping`: %w", err)
	}

	serviceProcess, err := service.New(pl.Services)
	if err != nil {
		return p, fmt.Errorf("error in `service`: %w", err)
	}

	silenceProcess, err := silence.New(pl.Silences)
	if err != nil {
		return p, fmt.Errorf("error in `silence`: %w", err)
	}

	ratelimitProcess, err := ratelimit.New(pl.Ratelimits)
	if err != nil {
		return p, fmt.Errorf("error in `ratelimit`: %w", err)
	}

	snoozeProcess, err := snooze.New()
	if err != nil {
		return p, fmt.Errorf("error in `snooze`: %w", err)
	}

	problemProcess, err := problem.New()
	if err != nil {
		return p, fmt.Errorf("error in `problem`: %w", err)
	}

	notificationProcess, err := notification.New(pl.Notifications, pl.DefaultDestinations)
	if err != nil {
		return p, fmt.Errorf("error in `notification`: %w", err)
	}

	storeProcess, err := store.New()
	if err != nil {
		return p, fmt.Errorf("error in `store`: %w", err)
	}

	p.processes = []SubProcessor{
		&timestamp.Processor{},
		transformProcess,
		profileProcess,
		groupingProcess,
		serviceProcess,
		silenceProcess,
		ratelimitProcess,
		snoozeProcess,
		&activecheck.Processor{},
		problemProcess,
		notificationProcess,
		storeProcess,
	}

	return p, nil
}

func (p *Processor) Process(ctx context.Context, item *models.Log) *decision.Decision {
	ctx, span := tracer.Start(ctx, "processing")
	defer span.End()

	// if item.TraceID != "" {
	// 	ctx = context.WithValue(ctx, tracing.ContextValue{}, item.TraceID)
	// }

	for _, subprocessor := range p.processes {
		if decision := subprocessor.Process(ctx, item); decision.Stop {
			return decision
		}
	}

	return decision.OK()
}
