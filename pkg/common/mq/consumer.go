package mq

import (
	"context"
	"sync"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

var onces = map[StreamName]*sync.Once{}
var consumers = map[StreamName]jetstream.Consumer{}

func Consumer(name StreamName) jetstream.Consumer {
	once, ok := onces[name]
	if !ok {
		once = &sync.Once{}
		onces[name] = once
	}
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cfg := jetstream.ConsumerConfig{}
		var err error
		consumers[name], err = streams[name].CreateOrUpdateConsumer(ctx, cfg)
		if err != nil {
			log.Fatal(err)
		}
	})
	return consumers[name]
}
