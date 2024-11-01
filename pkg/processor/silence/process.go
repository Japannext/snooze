package silence

import (
	"context"
	"fmt"

	"github.com/japannext/snooze/pkg/models"
)

func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracer.Start(ctx, "silence")
	defer span.End()

	for _, s := range silences {
		var match bool
		var err error
		// Condition
		if s.internal.condition != nil {
			match, err = s.internal.condition.MatchLog(ctx, item)
			if err != nil {
				log.Warnf("failed to match condition `%s` in silence `%s`: %s", s.Name, s.If, err)
				continue
			}
			if !match {
				continue
			}
		}

		// Schedule
		if s.Schedule != nil {
			match = s.Schedule.Match(&item.Timestamp.Actual.Time)
			if !match {
				continue
			}
		}

		item.Mute.Silence(fmt.Sprintf("Silenced by rule %s", s))
	}

	return nil
}
