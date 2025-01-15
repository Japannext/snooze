package processor

import (
	"context"
	"fmt"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/processor/decision"
	"github.com/japannext/snooze/pkg/processor/grouping"
	"github.com/japannext/snooze/pkg/processor/notification"
	"github.com/japannext/snooze/pkg/processor/profile"
	"github.com/japannext/snooze/pkg/processor/ratelimit"
	"github.com/japannext/snooze/pkg/processor/silence"
	"github.com/japannext/snooze/pkg/processor/snooze"
	"github.com/japannext/snooze/pkg/processor/store"
	"github.com/japannext/snooze/pkg/processor/transform"
)

type SubProcessor interface {
	//Name() string
	Process(ctx context.Context, item *models.Log) *decision.Decision
}

// An object that can process a log, and contain the static config
// of all processors
type Processor struct {
	transform,
	profile,
	grouping,
	silence,
	ratelimit,
	snooze,
	notification,
	store SubProcessor

	processes []SubProcessor
}

func NewProcessor(cfg *Config) (*Processor, error) {
	p := &Processor{}

	transformProcess, err := transform.New(cfg.Transform)
	if err != nil {
		return p, fmt.Errorf("error in `transform`: %w", err)
	}

	profileProcess, err := profile.New(cfg.Profile)
	if err != nil {
		return p, fmt.Errorf("error in `profile`: %w", err)
	}

	groupingProcess, err := grouping.New(cfg.Grouping)
	if err != nil {
		return p, fmt.Errorf("error in `grouping`: %w", err)
	}

	silenceProcess, err := silence.New(cfg.Silence)
	if err != nil {
		return p, fmt.Errorf("error in `silence`: %w", err)
	}

	ratelimitProcess, err := ratelimit.New(cfg.Ratelimit)
	if err != nil {
		return p, fmt.Errorf("error in `ratelimit`: %w", err)
	}

	snoozeProcess, err := snooze.New(cfg.Snooze)
	if err != nil {
		return p, fmt.Errorf("error in `snooze`: %w", err)
	}

	notificationProcess, err := notification.New(cfg.Notification)
	if err != nil {
		return p, fmt.Errorf("error in `notification`: %w", err)
	}

	storeProcess, err := store.New(cfg.Store)
	if err != nil {
		return p, fmt.Errorf("error in `store`: %w", err)
	}

	p.processes = []SubProcessor{
		transformProcess,
		profileProcess,
		groupingProcess,
		silenceProcess,
		ratelimitProcess,
		snoozeProcess,
		notificationProcess,
		storeProcess,
	}

	return p, nil
}

func (p *Processor) Process(ctx context.Context, item *models.Log) *decision.Decision {

	for _, subprocessor := range p.processes {
		if decision := subprocessor.Process(ctx, item); decision.Stop {
			return decision
		}
	}

	return decision.OK()
}
