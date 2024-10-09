package timestamp

import (
	"context"
	"sync"

	"github.com/japannext/snooze/pkg/models"
)

func Batch(ctx context.Context, items []*models.Log) {

	var wg sync.WaitGroup

	for _, item := range items {
		wg.Add(1)
		go func () {
			if item.Timestamp.Actual > item.Timestamp.Observed {
				item.Timestamp.Display = item.Timestamp.Observed
				item.Timestamp.Warning = "future"
			} else if item.Timestamp.Actual == 0 {
				item.Timestamp.Display = item.Timestamp.Observed
				item.Timestamp.Warning = "missing"
			} else {
				item.Timestamp.Display = item.Timestamp.Actual
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
