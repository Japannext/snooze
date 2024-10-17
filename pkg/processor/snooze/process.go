package snooze

import (
	"context"

	"github.com/japannext/snooze/pkg/models"
)

func Process(ctx context.Context, alert *models.Log) error {
	ctx, span := tracer.Start(ctx, "snooze")
	defer span.End()

	// TODO

	return nil
}
