package store

import (
	"context"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/processor/tracing"
)

func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracing.TRACER.Start(ctx, "store")
	defer span.End()
	if item.Mute.SkipStorage {
		return nil
	}
	_, err := opensearch.Store(ctx, models.LOG_INDEX, item)
	if err != nil {
		return err
	}
	storedLogs.Inc()
	log.Debugf("Successfully stored log %s", item)
	return nil
}
