package tracing

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
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

// Initialize the default tracer
func Init(serviceName string) {
	provider := NewTracerProvider(serviceName)
	otel.SetTracerProvider(provider)
}
