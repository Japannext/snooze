package otel

import (
	"fmt"
	"time"

	"github.com/japannext/snooze/pkg/models"
	commonv1 "go.opentelemetry.io/proto/otlp/common/v1"
	logv1 "go.opentelemetry.io/proto/otlp/logs/v1"
	resv1 "go.opentelemetry.io/proto/otlp/resource/v1"
)

const (
	SOURCE_KIND = "otel"
)

// Convert an opentelemetry format to the snooze native format.
func convertLog(resource *resv1.Resource, scope *commonv1.InstrumentationScope, lr *logv1.LogRecord) *models.Log {
	var item *models.Log

	item.Source = models.Source{Kind: SOURCE_KIND, Name: config.SourceName}

	// Timestamps
	item.ActualTime = models.Time{Time: time.Unix(0, int64(lr.GetTimeUnixNano()))}
	if lr.GetObservedTimeUnixNano() == 0 {
		item.ObservedTime = models.TimeNow()
	} else {
		item.ObservedTime = models.Time{Time: time.Unix(0, int64(lr.GetObservedTimeUnixNano()))}
	}

	item.SeverityText = lr.GetSeverityText()
	item.SeverityNumber = int32(lr.GetSeverityNumber())

	item.Labels = map[string]string{}
	for key, value := range kvToMap(resource.GetAttributes()) {
		item.Labels[fmt.Sprintf("otel.resource.%s", key)] = value
	}

	body := (&AnyValue{lr.GetBody()}).ToMap()
	for key, value := range body {
		if key == "message" {
			item.Message = value
		} else {
			item.Labels[fmt.Sprintf("otel.body.%s", key)] = value
		}
	}

	return item
}
