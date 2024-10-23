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

	if item.Mute.SkipStorage {
		log.Debugf("skipping storage")
		tracing.SetAttribute(span, "mute.skipStorage", "true")
		tracing.SetAttribute(span, "mute", item.Mute.Reason)
		return nil
	}
	tracing.SetAttribute(span, "mute.skipStorage", "false")

	return storeQ.PublishData(ctx, opensearch.Create(models.LOG_INDEX, item))
}
