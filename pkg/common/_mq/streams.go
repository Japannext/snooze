package mq

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

var streams = map[StreamName]jetstream.Stream{}

type StreamName string

const (
	NOTIFY_STREAM  StreamName = "NOTIFY"
	PROCESS_STREAM StreamName = "PROCESS"
	STORE_STREAM   StreamName = "STORE"
)

func SetupStreams() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	streams[PROCESS_STREAM], err = js.CreateStream(ctx, jetstream.StreamConfig{
		Name: string(PROCESS_STREAM),
		Retention: jetstream.WorkQueuePolicy,
		Subjects: []string{"PROCESS.logs"},
	})
	if err != nil {
		log.Fatal(err)
	}

	streams[NOTIFY_STREAM], err = js.CreateStream(ctx, jetstream.StreamConfig{
		Name: string(NOTIFY_STREAM),
		Retention: jetstream.WorkQueuePolicy,
		Subjects: []string{"NOTIFY.*"},
	})
	if err != nil {
		log.Fatal(err)
	}

	streams[STORE_STREAM], err = js.CreateStream(ctx, jetstream.StreamConfig{
		Name: string(STORE_STREAM),
		Retention: jetstream.WorkQueuePolicy,
		Subjects: []string{"STORE.*"},
	})
	if err != nil {
		log.Fatal(err)
	}
}
