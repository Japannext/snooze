package mq

import (
	"sync"

	"github.com/nats-io/nats.go/jetstream"
)

var (
	processSubOnce, processPubOnce, processStreamOnce sync.Once
	processStream jetstream.Stream
	processSub *Sub
	processPub *Pub
)

func getProcessStream() jetstream.Stream {
	processStreamOnce.Do(func() {
		client := getClient()
		processStream = client.setupStream(jetstream.StreamConfig{
			Name: "PROCESS",
			Retention: jetstream.WorkQueuePolicy,
			Subjects: []string{"PROCESS.logs"},
		})
	})
	return processStream
}

func ProcessSub() *Sub {
	processSubOnce.Do(func() {
		client := getClient()
		stream := getProcessStream()
		processSub = client.Consumer(stream, jetstream.ConsumerConfig{
			Name: "process",
			Durable: "process",
		})
	})
	return processSub
}

func ProcessPub() *Pub {
	processPubOnce.Do(func() {
		client := getClient()
		getProcessStream()
		processPub = client.Pub().WithSubject("PROCESS.logs")
	})
	return processPub
}
