package mq

import (
	"context"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/propagation"
)

var log = logrus.WithFields(logrus.Fields{"module": "mq"})

var (
	conn *nats.Conn
	js jetstream.JetStream
)

func SetupClient() {
	initConfig()
	conn, err := nats.Connect(config.URL)
	if err != nil {
		log.Fatal(err)
	}
	js, err = jetstream.New(conn)
	if err != nil {
		log.Fatal(err)
	}
}

type Client struct {
	conn *nats.Conn
	js jetstream.JetStream
}

func newClient() *Client {
	initConfig()
	conn, err := nats.Connect(config.URL)
	if err != nil {
		log.Fatal(err)
	}
	js, err = jetstream.New(conn)
	if err != nil {
		log.Fatal(err)
	}
	return &Client{conn: conn, js: js}
}

func (client *Client) Subscribe(subject string) *Subscription {
	return &Subscription{
		client: client,
		Subject: subject,
	}
}

type Subscription struct {
	client *Client
	consumer jetstream.Consumer
	StreamName string
	ConsumerName string
	Subject string
}

type MsgWithContext struct {
	Msg jetstream.Msg
	Context context.Context
}

func (msg MsgWithContext) Extract() (jetstream.Msg, context.Context) {
	return msg.Msg, msg.Context
}

var propagator = propagation.TraceContext{}

func (sub *Subscription) Fetch(size int) ([]MsgWithContext, error) {
	batch, err := sub.consumer.Fetch(size)
	if err != nil {
		return []MsgWithContext{}, err
	}

	var msgs = []MsgWithContext{}
	for msg := range batch.Messages() {
		ctx := context.Background()

		// Opentelemetry trace context
		propagator.Extract(ctx, propagation.HeaderCarrier(msg.Headers()))

		// In-queue custom metric
		observeDelay(sub.StreamName, msg)

		msgs = append(msgs, MsgWithContext{Msg: msg, Context: ctx})
	}
	return msgs, nil
}

const X_SNOOZE_PUBLISHED_TIME = "X-Snooze-Published-Time"

func observeDelay(streamName string, msg jetstream.Msg) {
	values, ok := msg.Headers()[X_SNOOZE_PUBLISHED_TIME]
	if !ok {
		return
	}
	if len(values) == 0 {
		return
	}
	i, err := strconv.Atoi(values[0])
	if err != nil {
		return
	}
	start := time.UnixMilli(int64(i))
	delay := time.Since(start)
	inQueue.WithLabelValues(streamName).Observe(float64(delay.Seconds()))
}

func Init() {
	SetupClient()
	SetupStreams()
}
