package profile

import (
	"context"
	"fmt"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/processor/decision"
	"go.opentelemetry.io/otel"
)

type Processor struct {
	mapper *FastMapper
}

type Config struct {
	Profiles []Profile `json:"profiles" yaml:"profiles"`
}

func New(cfg Config) (*Processor, error) {
	p := &Processor{}

	for _, prf := range cfg.Profiles {
		if err := prf.Load(); err != nil {
			return p, fmt.Errorf("failed to load profile '%s': %w", prf.Name, err)
		}
	}

	m, err := NewFastMapper(cfg.Profiles)
	if err != nil {
		return p, fmt.Errorf("failed to initialize fast mapper: %w", err)
	}
	p.mapper = m

	return p, nil
}

type Profile struct {
    // Name of the profile group
    Name string `json:"name" yaml:"name"`
    // The main condition for a log to match this rule. Used to
    // reduce the amount of processing by the use of maps.
    // Examples: process=sshd, service.name=keycloak, k8s.statefulset.name=postgresql
    Switch Kv `json:"switch" yaml:"switch"`
    // Patterns and actions to apply to logs matching this pattern
    Patterns []*Pattern `json:"patterns" yaml:"patterns"`
}

func (prf *Profile) Load() error {
    for _, pattern := range prf.Patterns {
        if err := pattern.Load(); err != nil {
            return fmt.Errorf("in profile=%s, pattern=%s: %s", prf.Name, pattern.Name, err)
        }
    }

	return nil
}

func (p *Processor) Process(ctx context.Context, item *models.Log) *decision.Decision {
	ctx, span := otel.Tracer("snooze").Start(ctx, "profile")
	defer span.End()

	for _, prf := range p.mapper.GetMatches(item) {
		for _, pattern := range prf.Patterns {
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
