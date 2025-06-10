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
	kv *kv.Kv[Profile]
}

func New(cfgs []pipeline.Profile) (*Processor, error) {

	kv := kv.New[Profile]()
	for _, cfg := range cfgs {
		p, err := NewProfile(cfg)
		if err != nil {
			return &Processor{}, fmt.Errorf("failed to initialize profile %s: %w", cfg.Name, err)
		}
		if err := kv.Add(p.Selector, &p); err != nil {
			return &Processor{}, fmt.Errorf("error in profile %s: %w", cfg.Name, err)
		}
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

	profiles := []*Profile{}

	// Find all profiles from the selectors.
	for _, h := range item.SelectorHashes {
		profiles = append(profiles, p.kv.AllMatches(h)...)
	}

	for _, prf := range profiles {
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
