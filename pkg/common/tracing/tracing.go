package tracing

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

var once sync.Once
var tracer trace.Tracer

func Tracer(name string) trace.Tracer {
	once.Do(func() {
		ctx := context.Background()
		exporter, err := otlptracegrpc.New(ctx)
		if err != nil {
			log.Fatal(err)
		}
		res := resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
		)
		provider := tracesdk.NewTracerProvider(
			tracesdk.WithBatcher(exporter, tracesdk.WithBatchTimeout(time.Second)),
			tracesdk.WithResource(res),
		)
		otel.SetTracerProvider(provider)
		tracer = otel.Tracer(name)
	})

	return tracer
}
