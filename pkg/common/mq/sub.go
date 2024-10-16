package mq

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/propagation"
)

func (client *Client) Subscribe(stream jetstream.Stream, name string) *Sub {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name: name,
		Durable: name,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Sub{
		client: client,
		consumer: consumer,
	}
}

type Sub struct {
	client *Client
	consumer jetstream.Consumer
}

type MsgWithContext struct {
	Msg jetstream.Msg
	Context context.Context
}

func (msg MsgWithContext) Extract() (jetstream.Msg, context.Context) {
	return msg.Msg, msg.Context
}

var propagator = propagation.TraceContext{}

func (sub *Sub) Fetch(size int, opts... jetstream.FetchOpt) ([]MsgWithContext, error) {
	batch, err := sub.consumer.Fetch(size, opts...)
	if err != nil {
		return []MsgWithContext{}, err
	}

	var msgs = []MsgWithContext{}
	for msg := range batch.Messages() {
		ctx := context.Background()

		// Opentelemetry trace context
		propagator.Extract(ctx, propagation.HeaderCarrier(msg.Headers()))

		// In-queue custom metric
		observeDelay("", msg)

		msgs = append(msgs, MsgWithContext{Msg: msg, Context: ctx})
	}
	return msgs, nil
}

