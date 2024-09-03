package profile

import (
	"context"
	// "github.com/sirupsen/logrus"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func Process(ctx context.Context, item *api.Log) error {
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
