package transform

import (
	"context"
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func Process(ctx context.Context, item *api.Log) error {
	for _, tr := range transforms {
		if tr.internal.condition != nil {
			match, err := tr.internal.condition.MatchLog(ctx, item)
			if err != nil {
				return err
			}
			if !match {
				continue
			}
		}
		if err := tr.internal.transform.Process(item); err != nil {
			return err
		}
	}

	return nil
}
