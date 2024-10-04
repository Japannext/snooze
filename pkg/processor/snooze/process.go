package snooze

import (
	"context"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/processor/tracing"
)

func Process(ctx context.Context, alert *api.Log) error {
	ctx, span := tracing.TRACER.Start(ctx, "snooze")
	defer span.End()

	// TODO

	return nil
}
