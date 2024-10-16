package mq

import (
	"context"
	"encoding/json"
	"time"
	"strconv"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func newHeader() nats.Header {
	return nats.Header(map[string][]string{
		X_SNOOZE_PUBLISHED_TIME: []string{strconv.Itoa(int(time.Now().UnixMilli()))},
	})
}

func PublishAsync(subject string, item interface{}) (jetstream.PubAckFuture, error) {
	data, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	msg := &nats.Msg{Subject: subject, Data: data, Header: newHeader()}
	return js.PublishMsgAsync(msg)
}

func Publish(ctx context.Context, subject string, item interface{}) error {
	body, err := json.Marshal(item)
	if err != nil {
		return err
	}
	msg := &nats.Msg{Subject: subject, Data: body, Header: newHeader()}
	if _, err := js.PublishMsg(ctx, msg); err != nil {
		return err
	}
	return nil
}

func PublishAsyncComplete() (<- chan struct{}) {
	return js.PublishAsyncComplete()
}
