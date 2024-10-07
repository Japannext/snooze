package snooze

import (
	"context"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/processor/tracing"
)

func Process(ctx context.Context, alert *models.Log) error {
	ctx, span := tracing.TRACER.Start(ctx, "snooze")
	defer span.End()

	// TODO

	return nil
}
