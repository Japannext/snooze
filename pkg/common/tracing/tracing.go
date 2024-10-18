package tracing

import (
	"context"
	"fmt"
	"time"
	"strconv"

	log "github.com/sirupsen/logrus"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

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

func SetAttribute(span trace.Span, key, value string) {
	if value == "" {
		return
	}
	span.SetAttributes(attribute.KeyValue{
		Key: attribute.Key(key),
		Value: attribute.StringValue(value),
	})
}

func SetMap(span trace.Span, prefix string, kv map[string]string) {
	for key, value := range kv {
		SetAttribute(span, fmt.Sprintf("%s.%s", prefix, key), value)
	}
}

func SetLog(span trace.Span, item *models.Log) {
	SetAttribute(span, "log.source.kind", item.Source.Kind)
	SetAttribute(span, "log.source.name", item.Source.Name)
	SetMap(span, "log.identity", item.Identity)
	SetMap(span, "log.labels", item.Labels)
	SetAttribute(span, "log.severityText", item.SeverityText)
	SetAttribute(span, "log.severityNumber", strconv.Itoa(int(item.SeverityNumber)))
	SetAttribute(span, "log.message", item.Message)
}

// Initialize the default tracer
func Init(serviceName string) {
	provider := NewTracerProvider(serviceName)
	otel.SetTracerProvider(provider)
}
