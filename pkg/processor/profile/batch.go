package profile

import (
	"context"
	"sync"

	"github.com/japannext/snooze/pkg/processor/tracing"
	"github.com/japannext/snooze/pkg/models"
)

func Batch(ctx context.Context, items []*models.Log) error {
	ctx, span := tracing.TRACER.Start(ctx, "profile")
	defer span.End()
	var wg sync.WaitGroup
	for _, item := range items {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Process(ctx, item)
		}()
	}
	wg.Wait()

	return nil
}
