package transform

import (
	"context"

	"github.com/japannext/snooze/pkg/common/tracing"
	"github.com/japannext/snooze/pkg/models"
)

var tracer = tracing.Tracer("snooze-process")

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
		if err := tr.internal.transform.Process(item); err != nil {
			return err
		}
	}

	return nil
}
