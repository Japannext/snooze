package store

import (
	"context"

	"github.com/japannext/snooze/pkg/models"
)

func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracer.Start(ctx, "store")
	defer span.End()

	if item.Mute.SkipStorage {
		return nil
	}

	return storeQ.Publish(ctx, item)
}
