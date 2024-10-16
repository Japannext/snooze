package mq

import (
	"context"
	"sync"
	"time"
	"strconv"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var onces = map[StreamName]*sync.Once{}
var consumers = map[string]jetstream.Consumer{}

func Consumer(stream StreamName, name string) jetstream.Consumer {
	once, ok := onces[stream]
	if !ok {
		once = &sync.Once{}
		onces[stream] = once
	}
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cfg := jetstream.ConsumerConfig{
			Name: name,
			Durable: name,
		}
		var err error
		consumers[name], err = streams[stream].CreateOrUpdateConsumer(ctx, cfg)
		if err != nil {
			log.Fatal(err)
		}
	})
	return consumers[name]
}

// Adding observability to the consumer (traces and time-in-queue metrics)
type Consumer2 struct {
	js jetstream.Consumer
	streamName string
}

func NewConsumer(stream jetstream.Stream, streamName, name string) *Consumer2 {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := jetstream.ConsumerConfig{
		Name: name,
		Durable: name,
	}
	var err error
	consumer, err := stream.CreateOrUpdateConsumer(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &Consumer2{
		js: consumer,
	}
}

func (c *Consumer2) Fetch(batchSize int, opts ...jetstream.FetchOpt) ([]jetstream.Msg, error) {

	batch, err := c.js.Fetch(batchSize, opts...)
	if err != nil {
		return []jetstream.Msg{}, err
	}

	var msgs = []jetstream.Msg{}
	for msg := range batch.Messages() {
		msgs = append(msgs, msg)
		observeDelay(c, msg)
	}

	return msgs, nil
}
