package profile

import (
	"context"
	// "github.com/sirupsen/logrus"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/processor/tracing"
)

func Process(ctx context.Context, item *api.Log) error {
	ctx, span := tracing.TRACER.Start(ctx, "profile")
	defer span.End()
	for _, rule := range fastMapper.GetMatches(item) {
		reject := rule.Process(ctx, item)
		if reject {
			return nil
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
