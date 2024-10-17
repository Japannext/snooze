package mq

import (
	"context"
	"encoding/json"
	"maps"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/propagation"
)

func (client *Client) Pub() *Pub {
	return &Pub{client: client, subject: "", headers: make(map[string]string)}
}

type Pub struct {
	client *Client
	subject string
	headers map[string]string
}

func (pub *Pub) WithSubject(subject string) *Pub {
	return &Pub{
		client: pub.client,
		headers: pub.headers,
		subject: subject,
	}
}

func (pub *Pub) WithHeader(key, value string) *Pub {
	headers := map[string]string{}
	maps.Copy(pub.headers, headers)
	headers[key] = value
	return &Pub{
		client: pub.client,
		subject: pub.subject,
		headers: headers,
	}
}

func (pub *Pub) WithIndex(index string) *Pub {
	return pub.WithHeader(X_SNOOZE_STORE_INDEX, index)
}

func (pub *Pub) Publish(ctx context.Context, item interface{}) error {
	ctx, span := tracer.Start(ctx, "Publish")
	defer span.End()
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	header := make(nats.Header)
	for k, v := range pub.headers {
		header.Add(k, v)
	}
	// Metrics
	injectPublishTime(&header)
	// Opentelemetry
	propagator := propagation.TraceContext{}
	propagator.Inject(ctx, propagation.HeaderCarrier(header))
	msg := &nats.Msg{Subject: pub.subject, Data: data, Header: header}

	if _, err := pub.client.js.PublishMsg(ctx, msg); err != nil {
		return err
	}
	return nil
}

func (pub *Pub) Wait() <-chan struct{} {
	return pub.client.js.PublishAsyncComplete()
}
