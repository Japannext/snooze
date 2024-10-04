package transform

import (
	"context"
	"sync"

	"github.com/japannext/snooze/pkg/processor/tracing"
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func Batch(ctx context.Context, items []*api.Log) error {
	ctx, span := tracing.TRACER.Start(ctx, "transform")
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
