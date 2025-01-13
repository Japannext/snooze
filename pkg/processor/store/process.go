package store

import (
	"context"

	"github.com/japannext/snooze/pkg/common/opensearch/format"
	"github.com/japannext/snooze/pkg/common/tracing"
	"github.com/japannext/snooze/pkg/models"
)

func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracer.Start(ctx, "store")
	defer span.End()

	tracing.SetBool(span, "status.skipStorage", item.Status.SkipStorage)
	if item.Status.SkipStorage {
		log.Debugf("skipping storage")
		return nil
	}

	return storeQ.PublishData(ctx, &format.Create{
		Index: models.LogIndex,
		Item:  item,
	})
}
