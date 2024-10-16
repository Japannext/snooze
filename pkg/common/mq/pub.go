package mq

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
)

func (client *Client) Pub(defaultSubject string) *Pub {
	return &Pub{client: client, defaultSubject: defaultSubject}
}

type Pub struct {
	client *Client
	defaultSubject string
}

func (pub *Pub) PublishWithSubject(ctx context.Context, subject string, item interface{}) error {
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	msg := &nats.Msg{Subject: subject, Data: data}
	injectPublishTime(msg)
	if _, err := pub.client.js.PublishMsg(ctx, msg); err != nil {
		return err
	}
	return nil
}

func (pub *Pub) Publish(ctx context.Context, item interface{}) error {
	return pub.PublishWithSubject(ctx, pub.defaultSubject, item)
}

func (pub *Pub) Wait() <-chan struct{} {
	return pub.client.js.PublishAsyncComplete()
}
