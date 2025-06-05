package profile

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
	kv *kv.Kv[*models.Log, Profile]
}

func New(cfgs []pipeline.Profile) (*Processor, error) {

	profiles := []Profile{}
	for _, cfg := range cfgs {
		p, err := NewProfile(cfg)
		if err != nil {
			return &Processor{}, fmt.Errorf("failed to initialize profile %s: %w", cfg.Name, err)
		}
		profiles = append(profiles, p)
	}

	kv, err := kv.New[*models.Log](profiles)
	if err != nil {
		return &Processor{}, fmt.Errorf("failed to initialize fast mapper: %w", err)
	}

	return &Processor{
		kv: kv,
	}, nil
}

type Profile struct {
	pipeline.Profile

	patterns []Pattern
}

func NewProfile(cfg pipeline.Profile) (Profile, error) {
	patterns := []Pattern{}
	for _, patternConfig := range cfg.Patterns {
		pattern, err := NewPattern(patternConfig)
		if err != nil {
			return Profile{}, fmt.Errorf("in profile=%s, pattern=%s: %s", cfg.Name, patternConfig.Name, err)
		}
		patterns = append(patterns, pattern)
	}

	return Profile{
		Profile:  cfg,
		patterns: patterns,
	}, nil
}

func (p *Processor) Process(ctx context.Context, item *models.Log) *decision.Decision {
	ctx, span := otel.Tracer("snooze").Start(ctx, "profile")
	defer span.End()

	for _, prf := range p.kv.AllMatches(item) {
		for _, pattern := range prf.patterns {
			match, reject := pattern.Process(ctx, item)
			if reject {
				return decision.OK()
			}

			if match {
				item.Profile = prf.Name
				item.Pattern = pattern.Name

				return decision.OK()
			}
		}
	}

	return decision.OK()
}
