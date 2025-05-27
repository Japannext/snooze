package silence

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"

	"github.com/japannext/snooze/components/processor/decision"
	"github.com/japannext/snooze/components/processor/metrics"
	"github.com/japannext/snooze/internal/lang"
	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/pipeline"
	"github.com/japannext/snooze/internal/schedule"
	"github.com/japannext/snooze/internal/tracing"
)

type Processor struct {
	silences []Silence
}

func New(cfgs []pipeline.Silence) (*Processor, error) {
	p := &Processor{}

	for _, cfg := range cfgs {
		s, err := NewSilence(cfg)
		if err != nil {
			return &Processor{}, fmt.Errorf("in silence '%s': %w", s.Name, err)
		}
		p.silences = append(p.silences, s)
	}

	return p, nil
}

type Silence struct {
	pipeline.Silence
	condition *lang.Condition
}

func NewSilence(cfg pipeline.Silence) (Silence, error) {
	var err error

	s := Silence{Silence: cfg}

	s.condition, err = lang.NewCondition(cfg.If)
	if err != nil {
		return Silence{}, fmt.Errorf("error in condition `%s`: %w", cfg.If, err)
	}

	if cfg.Schedule == nil {
		s.Schedule = schedule.Default()
	}

	if err := s.Schedule.Load(); err != nil {
		return Silence{}, fmt.Errorf("error in schedule: %w", err)
	}

	return s, nil
}

func (p *Processor) Process(ctx context.Context, item *models.Log) *decision.Decision {
	ctx, span := otel.Tracer("snooze").Start(ctx, "silence")
	defer span.End()

	for _, s := range p.silences {
		var match bool
		var err error

		// Condition
		if s.condition != nil {
			match, err = s.condition.MatchLog(ctx, item)
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
