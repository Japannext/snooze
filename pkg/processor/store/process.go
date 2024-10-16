package store

import (
	"context"

	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/models"
)

var storeQ = mq.StorePub("STORE.v2-logs")

func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracer.Start(ctx, "store")
	defer span.End()
	return storeQ.Publish(ctx, item)
}
