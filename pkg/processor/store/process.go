package store

import (
	"context"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/tracing"
)

func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracer.Start(ctx, "store")
	defer span.End()

	if item.Status.SkipStorage {
		log.Debugf("skipping storage")
		tracing.SetAttribute(span, "mute.skipStorage", "true")
		return nil
	}
	tracing.SetAttribute(span, "mute.skipStorage", "false")

	return storeQ.PublishData(ctx, opensearch.Create(models.LOG_INDEX, item))
}
