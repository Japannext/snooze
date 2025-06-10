package mq

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/japannext/snooze/internal/tracing"
)

func (client *Client) Consumer(stream jetstream.Stream, cfg jetstream.ConsumerConfig, opts ...jetstream.FetchOpt) *Sub {
	log.Debugf("subscribing consumer '%s'", cfg.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	consumer, err := stream.CreateOrUpdateConsumer(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to create/update consumer '%s': %s", cfg.Name, err)
	}

	streamName := stream.CachedInfo().Config.Name

	if len(opts) == 0 {
		opts = append(opts, jetstream.FetchMaxWait(1*time.Second))
	}

	return &Sub{
		client:     client,
		consumer:   consumer,
		streamName: streamName,
		fetchOpts:  opts,
	}
}

type Sub struct {
	client     *Client
	consumer   jetstream.Consumer
	streamName string
	fetchOpts  []jetstream.FetchOpt
}

type MsgWithContext struct {
	Msg     jetstream.Msg
	Context context.Context
}

func (msg MsgWithContext) Extract() (jetstream.Msg, context.Context) {
	return msg.Msg, msg.Context
}

func (sub *Sub) Fetch(size int) ([]MsgWithContext, error) {
	//ctx, fetchSpan := tracer.Start(context.Background(), "Fetch")
	// defer fetchSpan.End()

	batch, err := sub.consumer.Fetch(size, sub.fetchOpts...)
	if err != nil {
		return nil, err
	}

	msgs := []MsgWithContext{}

	for msg := range batch.Messages() {
		// Opentelemetry trace context
		// headers := msg.Headers()
		// propagator := propagation.TraceContext{}
		prop := otel.GetTextMapPropagator()
		headers := propagation.HeaderCarrier(msg.Headers())
		//msgCtx := propagator.Extract(context.Background(), propagation.HeaderCarrier(headers))
		msgCtx := prop.Extract(context.Background(), headers)
		log.Debugf("traceID (propagated): %s", getTraceID(msgCtx))

		// Opentelemetry msg-level span
		msgCtx, span := tracer.Start(msgCtx, "Fetch")
		for key := range headers {
			tracing.SetString(span, fmt.Sprintf("nats.header.%s", key), headers.Get(key))
		}
		tracing.SetString(span, "nats.data", string(msg.Data()))

		// In-queue custom metric
		start, ok := getPublishedTime(msg)
		if ok {
			observeDelay(sub.streamName, start)
		}

		msgs = append(msgs, MsgWithContext{Msg: msg, Context: msgCtx})
		span.End()
	}
	return msgs, nil
}

func getTraceID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		traceID := spanCtx.TraceID()
		return traceID.String()
	}
	return ""
}
