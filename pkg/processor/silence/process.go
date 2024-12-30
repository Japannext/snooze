package silence

import (
	"context"
	"fmt"

	"github.com/japannext/snooze/pkg/common/tracing"
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

		if ok := item.Status.Change(models.LogSilenced); ok {
			item.Status.Reason = fmt.Sprintf("Silenced by '%s'", s.Name)
			item.Status.SkipNotification = true
			silencedLogs.Inc()
		}
		if s.Drop {
			if ok := item.Status.Change(models.LogDropped); ok {
				item.Status.Reason = fmt.Sprintf("Dropped by silence '%s'", s.Name)
				item.Status.SkipNotification = true
				item.Status.SkipStorage = true
			}
		}
	}

	return nil
}
