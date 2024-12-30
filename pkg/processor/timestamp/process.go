package timestamp

import (
	"context"

	"github.com/japannext/snooze/pkg/models"
)

func Process(ctx context.Context, item *models.Log) {
	actual, observed := item.ActualTime.Time, item.ObservedTime.Time
	if actual.After(observed) {
		item.DisplayTime = item.ObservedTime
		// item.WarningTime = "future"
	} else if actual.IsZero() {
		item.DisplayTime = item.ObservedTime
		// item.Timestamp.Warning = "missing"
	} else {
		item.DisplayTime = item.ActualTime
	}
}
