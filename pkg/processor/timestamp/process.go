package timestamp

import (
	"context"

	"github.com/japannext/snooze/pkg/models"
)

func Process(ctx context.Context, item *models.Log) {
	actual, observed := item.Timestamp.Actual.Time, item.Timestamp.Observed.Time
	if actual.After(observed) {
		item.Timestamp.Display = item.Timestamp.Observed
		item.Timestamp.Warning = "future"
	} else if actual.IsZero() {
		item.Timestamp.Display = item.Timestamp.Observed
		item.Timestamp.Warning = "missing"
	} else {
		item.Timestamp.Display = item.Timestamp.Actual
	}
}
