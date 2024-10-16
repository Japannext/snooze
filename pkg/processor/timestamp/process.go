package timestamp

import (
	"context"

	"github.com/japannext/snooze/pkg/models"
)

func Process(ctx context.Context, item *models.Log) {
	if item.Timestamp.Actual > item.Timestamp.Observed {
		item.Timestamp.Display = item.Timestamp.Observed
		item.Timestamp.Warning = "future"
	} else if item.Timestamp.Actual == 0 {
		item.Timestamp.Display = item.Timestamp.Observed
		item.Timestamp.Warning = "missing"
	} else {
		item.Timestamp.Display = item.Timestamp.Actual
	}
}
