package transform

import (
	"context"

	"github.com/japannext/snooze/pkg/models"
)

func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracer.Start(ctx, "transform")
	defer span.End()
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

		ctx = context.WithValue(ctx, "capture", map[string]string{})
		for _, action := range tr.internal.actions {
			var err error
			ctx, err = action.Process(ctx, item)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
