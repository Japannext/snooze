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
		v, err := s.internal.condition.MatchLog(ctx, item)
		if err != nil {
			return err
		}
		if v {
			item.Mute.Silence(fmt.Sprintf("Silenced by rule %s", s))
			break
		}
	}

	return nil
}
