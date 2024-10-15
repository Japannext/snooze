package mq

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func PublishAsync(subject string, item interface{}) (jetstream.PubAckFuture, error) {
	body, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return js.PublishAsync(subject, body)
}

func Publish(ctx context.Context, subject string, item interface{}) error {
	body, err := json.Marshal(item)
	if err != nil {
		return err
	}
	msg := &nats.Msg{
		Subject: subject,
		Data: body,
	}
	if _, err := js.PublishMsg(ctx, msg); err != nil {
		return err
	}
	return nil
}
