package mq

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

func (client *Client) setupStream(cfg jetstream.StreamConfig) jetstream.Stream {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.js.CreateStream(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	return stream
}
