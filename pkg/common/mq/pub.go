package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/japannext/snooze/pkg/common/tracing"
)

func (client *Client) Pub() *Pub {
	return &Pub{client: client, subject: "", headers: make(map[string]string)}
}

type Pub struct {
	client  *Client
	subject string
	headers map[string]string
}

func (pub *Pub) WithSubject(subject string) *Pub {
	return &Pub{
		client:  pub.client,
		headers: pub.headers,
		subject: subject,
	}
}

const X_SNOOZE_WRITER_ACTION = "X-Snooze-Writer-Action"

func (pub *Pub) WithHeader(key, value string) *Pub {
	headers := map[string]string{}
	maps.Copy(pub.headers, headers)
	headers[key] = value
	return &Pub{
		client:  pub.client,
		subject: pub.subject,
		headers: headers,
	}
}

func (pub *Pub) WithIndex(index string) *Pub {
	return pub.WithHeader(X_SNOOZE_STORE_INDEX, index)
}

func (pub *Pub) WithAction(action string) *Pub {
	return pub.WithHeader(X_SNOOZE_WRITER_ACTION, action)
}

func (pub *Pub) publish(ctx context.Context, span trace.Span, data []byte) error {
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

	for key, _ := range header {
		tracing.SetString(span, fmt.Sprintf("nats.header.%s", key), header.Get(key))
	}
	tracing.SetString(span, "nats.data", string(data))

	if _, err := pub.client.js.PublishMsg(ctx, msg); err != nil {
		tracing.Error(span, err)
		return err
	}
	return nil
}

func (pub *Pub) PublishData(ctx context.Context, item Serializable) error {
	ctx, span := tracer.Start(ctx, "PublishData")
	defer span.End()

	data, err := item.Serialize()
	if err != nil {
		tracing.Error(span, err)
		return err
	}

	return pub.publish(ctx, span, data)
}

type Serializable interface {
	Serialize() ([]byte, error)
}

func (pub *Pub) Publish(ctx context.Context, item interface{}) error {
	ctx, span := tracer.Start(ctx, "Publish")
	defer span.End()

	data, err := json.Marshal(item)
	if err != nil {
		tracing.Error(span, err)
		return err
	}

	return pub.publish(ctx, span, data)
}

func (pub *Pub) Wait() <-chan struct{} {
	return pub.client.js.PublishAsyncComplete()
}
