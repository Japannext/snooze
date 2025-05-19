package timestamp

import (
	"context"

	"github.com/japannext/snooze/components/processor/decision"
	"github.com/japannext/snooze/internal/models"
)

type Processor struct{}

func (p *Processor) Process(_ context.Context, item *models.Log) *decision.Decision {
	actual, observed := item.ActualTime.Time, item.ObservedTime.Time

	switch {
	case actual.After(observed):
		item.DisplayTime = item.ObservedTime
	case actual.IsZero():
		item.DisplayTime = item.ObservedTime
	default:
		item.DisplayTime = item.ActualTime
	}

	return decision.OK()
}
