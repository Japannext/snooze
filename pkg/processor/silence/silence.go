package silence

import (
	"context"
	"fmt"

	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/common/schedule"
	"github.com/japannext/snooze/pkg/common/tracing"
	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/processor/decision"
	"github.com/japannext/snooze/pkg/processor/metrics"
	"go.opentelemetry.io/otel"
	log "github.com/sirupsen/logrus"
)

type Processor struct {
	silences []Silence
}

type Config struct {
	Silences []Silence `json:"silences" yaml:"silences"`
}

func New(cfg Config) (*Processor, error) {
	p := &Processor{}

	for _, sil := range cfg.Silences {
		if err := sil.Load(); err != nil {
			return p, fmt.Errorf("in silence '%s': %w", sil.Name, err)
		}
		p.silences = append(p.silences, sil)
	}

	return p, nil
}

type Silence struct {
	Name        string             `json:"name"                  yaml:"name"`
	Description string             `json:"description,omitempty" yaml:"description"`
	If          string             `json:"if"                    yaml:"if"`
	Schedule    *schedule.Schedule `yaml:",inline"`
	Drop        bool               `json:"drop"                  yaml:"drop"`

	internal struct {
		condition *lang.Condition
	}
}

func (s *Silence) Load() error {
	var err error

	s.internal.condition, err = lang.NewCondition(s.If)
	if err != nil {
		return fmt.Errorf("error in condition `%s`: %w", s.If, err)
	}

	if s.Schedule == nil {
		s.Schedule = schedule.Default()
	}

	if err := s.Schedule.Load(); err != nil {
		return fmt.Errorf("error in schedule: %w", err)
	}

	return nil
}

func (p *Processor) Process(ctx context.Context, item *models.Log) *decision.Decision {
	ctx, span := otel.Tracer("snooze").Start(ctx, "silence")
	defer span.End()

	for _, s := range p.silences {
		var match bool
		var err error

		// Condition
		if s.internal.condition != nil {
			match, err = s.internal.condition.MatchLog(ctx, item)
			if err != nil {
				log.Warnf("error while matching `%s` : %s", s.If, err)

				continue
			}

			tracing.SetBool(span, fmt.Sprintf("silence.%s.match", s.Name), match)

			if !match {
				continue
			}
		}

		// Schedule
		if s.Schedule != nil {
			match = s.Schedule.Match(&item.ActualTime.Time)
			tracing.SetBool(span, fmt.Sprintf("silence.%s.schedule", s.Name), match)

			if !match {
				continue
			}
		}

		if s.Drop {
			if ok := item.Status.Change(models.LogDropped); ok {
				item.Status.Reason = fmt.Sprintf("Silenced by '%s'", s.Name)
				item.Status.SkipNotification = true
				item.Status.SkipStorage = true

				metrics.SilencedLogs.Inc()
			}

			return decision.Done()
		}

		if ok := item.Status.Change(models.LogSilenced); ok {
			item.Status.Reason = fmt.Sprintf("Silenced by '%s'", s.Name)
			item.Status.SkipNotification = true

			metrics.SilencedLogs.Inc()
		}
	}

	return decision.OK()
}
