package otel

import (
	"time"
	"fmt"

	commonv1 "go.opentelemetry.io/proto/otlp/common/v1"
	logv1 "go.opentelemetry.io/proto/otlp/logs/v1"
	resv1 "go.opentelemetry.io/proto/otlp/resource/v1"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

const (
	SOURCE_KIND = "otel"
)

// Convert an opentelemetry format to the snooze native format
func convertLog(resource *resv1.Resource, scope *commonv1.InstrumentationScope, lr *logv1.LogRecord) *api.Log {
	var item *api.Log

	item.Source = api.Source{Kind: SOURCE_KIND, Name: config.SourceName}

	// Timestamps
	item.TimestampMillis = lr.TimeUnixNano / 1000 / 1000
	if lr.ObservedTimeUnixNano == 0 {
		item.ObservedTimestampMillis = uint64(time.Now().UnixMilli())
	} else {
		item.ObservedTimestampMillis = lr.ObservedTimeUnixNano / 1000 / 1000
	}

	item.SeverityText = lr.SeverityText
	item.SeverityNumber = int32(lr.SeverityNumber)

	item.Labels = map[string]string{}
	for key, value := range kvToMap(resource.Attributes) {
		item.Labels[fmt.Sprintf("otel.resource.%s", key)] = value
	}

	body := (&AnyValue{lr.Body}).ToMap()
	for key, value := range body {
		if key == "message" {
			item.Message = value
		} else {
			item.Labels[fmt.Sprintf("otel.body.%s", key)] = value
		}
	}

	return item
}
