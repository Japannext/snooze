package mq

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func (client *Client) Consumer(stream jetstream.Stream, cfg jetstream.ConsumerConfig) *Sub {
	log.Debugf("subscribing consumer '%s'", cfg.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	consumer, err := stream.CreateOrUpdateConsumer(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to create/update consumer '%s': %s", cfg.Name, err)
	}

	streamName := stream.CachedInfo().Config.Name

	return &Sub{
		client: client,
		consumer: consumer,
		streamName: streamName,
	}
}

type Sub struct {
	client *Client
	consumer jetstream.Consumer
	streamName string
}

type MsgWithContext struct {
	Msg jetstream.Msg
	Context context.Context
}

func (msg MsgWithContext) Extract() (jetstream.Msg, context.Context) {
	return msg.Msg, msg.Context
}

func (sub *Sub) Fetch(size int, opts... jetstream.FetchOpt) ([]MsgWithContext, error) {
	ctx, fetchSpan := tracer.Start(context.Background(), "Fetch")
	batch, err := sub.consumer.Fetch(size, opts...)
	defer fetchSpan.End()
	if err != nil {
		return nil, err
	}

	var msgs = []MsgWithContext{}
	for msg := range batch.Messages() {

		// Opentelemetry trace context
		propagator := propagation.TraceContext{}
		msgCtx := propagator.Extract(context.Background(), propagation.HeaderCarrier(msg.Headers()))

		// Opentelemetry msg-level span
		msgCtx, span := tracer.Start(msgCtx, "Fetch", trace.WithLinks(trace.LinkFromContext(ctx)))
		defer span.End()

		// In-queue custom metric
		start, ok := getPublishedTime(msg)
		if ok {
			observeDelay(sub.streamName, start)
		}

		msgs = append(msgs, MsgWithContext{Msg: msg, Context: msgCtx})
	}
	return msgs, nil
}

