package silence

import (
	"context"
	"fmt"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/processor/tracing"
)

func Process(ctx context.Context, item *api.Log) error {
	ctx, span := tracing.TRACER.Start(ctx, "silence")
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

// Process a batch of items
func Batch(ctx context.Context, items []*api.Log) {
	 for _, item := range items {
		Process(ctx, item)
	 }
}
