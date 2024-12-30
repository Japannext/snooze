package tracing

import (
	"context"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/japannext/snooze/pkg/models"
)

func NewTracerProvider(serviceName string) trace.TracerProvider {
	ctx := context.Background()
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
	)
	provider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter, tracesdk.WithBatchTimeout(time.Second)),
		tracesdk.WithResource(res),
	)
	return provider
}

func SetString(span trace.Span, key, value string) {
	span.SetAttributes(attribute.KeyValue{
		Key:   attribute.Key(key),
		Value: attribute.StringValue(value),
	})
}

func SetInt(span trace.Span, key string, value int) {
	span.SetAttributes(attribute.KeyValue{
		Key:   attribute.Key(key),
		Value: attribute.IntValue(value),
	})
}

func SetInt64(span trace.Span, key string, value int64) {
	span.SetAttributes(attribute.KeyValue{
		Key:   attribute.Key(key),
		Value: attribute.Int64Value(value),
	})
}

func SetFloat64(span trace.Span, key string, value float64) {
	span.SetAttributes(attribute.KeyValue{
		Key:   attribute.Key(key),
		Value: attribute.Float64Value(value),
	})
}

func SetBool(span trace.Span, key string, value bool) {
	span.SetAttributes(attribute.KeyValue{
		Key:   attribute.Key(key),
		Value: attribute.BoolValue(value),
	})
}

func SetTime(span trace.Span, key string, value time.Time) {
	span.SetAttributes(attribute.KeyValue{
		Key:   attribute.Key(key),
		Value: attribute.StringValue(value.Format(time.RFC3339)),
	})
}

func SetMap(span trace.Span, prefix string, kv map[string]string) {
	for key, value := range kv {
		SetString(span, fmt.Sprintf("%s.%s", prefix, key), value)
	}
}

func SetLog(span trace.Span, item *models.Log) {
	SetString(span, "log.source.kind", item.Source.Kind)
	SetString(span, "log.source.name", item.Source.Name)
	SetMap(span, "log.identity", item.Identity)
	SetMap(span, "log.labels", item.Labels)
	SetString(span, "log.severityText", item.SeverityText)
	SetInt(span, "log.severityNumber", int(item.SeverityNumber))
	SetString(span, "log.message", item.Message)
	SetTime(span, "log.actualTime", item.ActualTime.Time)
	SetTime(span, "log.observedTime", item.ObservedTime.Time)
	SetTime(span, "log.displayTime", item.DisplayTime.Time)
}

func HTTPFilter(req *http.Request) bool {
	if req.URL.Path == "/livez" ||
		req.URL.Path == "/readyz" ||
		req.URL.Path == "/metrics" {
		return false
	}
	return true
}

func GetTraceID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}
	return ""
}

func Warning(span trace.Span, err error) {
	span.RecordError(err)
}
func Error(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// Initialize the default tracer
func Init(serviceName string) {
	provider := NewTracerProvider(serviceName)
	otel.SetTracerProvider(provider)
}
