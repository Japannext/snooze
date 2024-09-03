package grouping

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/common/utils"
)

func Process(ctx context.Context, item *api.Log) error {
	if len(item.Group.Labels) != 0 { // existing grouping
		item.Group.Hash = utils.ComputeHash(item.Group.Labels)
		return nil
	}
	for _, group := range groupings {
		match, err := group.internal.condition.MatchLog(ctx, item)
		if err != nil {
			return err
		}
		if match {
			item.Group.Labels = make(map[string]string)
			for _, field := range group.internal.fields {
				value, err := lang.ExtractField(item, field)
				if err != nil {
					logrus.Warnf("Failed to match %s: %s", field, err)
					continue
				}
				item.Group.Labels[field.String()] = value
			}
			item.Group.Hash = utils.ComputeHash(item.Group.Labels)
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
