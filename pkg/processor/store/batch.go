package store

import (
	"context"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/processor/tracing"
)

func Batch(ctx context.Context, items []*api.Log) error {
	ctx, span := tracing.TRACER.Start(ctx, "store")
	defer span.End()

	// Removing logs that skip storage
	var logs []*api.Log
	for _, item := range items {
		if item.Mute.SkipStorage {
			continue
		}
		logs = append(logs, item)
	}

	if len(logs) == 0 {
		log.Debugf("Skipped %d logs", len(items))
		return nil
	}

	err := opensearch.StoreLogs(ctx, api.LOG_INDEX, logs)
	if err != nil {
		log.Warnf("failed to store batch: %s", err)
		return err
	}
	storedLogs.Add(float64(len(logs)))
	log.Debugf("Successfully stored %d logs", len(logs))

	return nil
}
