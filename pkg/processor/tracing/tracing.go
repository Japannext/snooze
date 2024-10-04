package tracing

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

var TRACER = otel.Tracer("snooze-process")

func Init() {
	ctx := context.Background()
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("snooze-process"),
	)
	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter, trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(provider)
}
