package otel

import (
  "fmt"

  commonv1 "go.opentelemetry.io/proto/otlp/common/v1"
  logv1 "go.opentelemetry.io/proto/otlp/logs/v1"
  resv1 "go.opentelemetry.io/proto/otlp/resource/v1"

  "github.com/japannext/snooze/common/api/v2"
)

// Convert an opentelemetry format to the snooze native format
func convertLog(resource *resv1.Resource, scope *commonv1.InstrumentationScope, lr *logv1.LogRecord) *v2.AlertEvent {
  var alert *v2.AlertEvent

  alert.Kind = "opentelemetry.io/logv1"

  // Timestamps
  alert.Timestamp = lr.TimeUnixNano
  if lr.ObservedTimeUnixNano == 0 {
    alert.ObservedTimestamp = timeNow()
  } else {
    alert.ObservedTimestamp = lr.ObservedTimeUnixNano
  }

  alert.SeverityText = lr.SeverityText
  alert.SeverityNumber = int32(lr.SeverityNumber)

  alert.Resource = kvToMap(resource.Attributes)
  alert.Attributes = kvToMap(lr.Attributes)

  alert.Body = lr.Body.GetStringValue()

  return alert
}
