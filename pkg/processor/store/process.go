package store

import (
	"context"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/processor/tracing"
)

func Process(ctx context.Context, item *api.Log) error {
	ctx, span := tracing.TRACER.Start(ctx, "store")
	defer span.End()
	if item.Mute.SkipStorage {
		return nil
	}
	_, err := opensearch.Store(ctx, api.LOG_INDEX, item)
	if err != nil {
		return err
	}
	storedLogs.Inc()
	log.Debugf("Successfully stored log %s", item)
	return nil
}

// Process a batch of items
func Batch(items []*api.Log) {
	// TODO
}
