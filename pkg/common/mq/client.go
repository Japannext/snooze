package mq

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/trace"
	"github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/common/tracing"
)

var log *logrus.Entry
var tracer trace.Tracer

type Client struct {
	conn *nats.Conn
	js jetstream.JetStream
}

func newClient() *Client {
	initConfig()
	conn, err := nats.Connect(config.URL)
	if err != nil {
		log.Fatalf("failed to connect to nats: %s", err)
	}
	js, err := jetstream.New(conn)
	if err != nil {
		log.Fatalf("failed to initialize jetstream: %s", err)
	}
	return &Client{conn: conn, js: js}
}

var client *Client

func Init() {
	log = logrus.WithFields(logrus.Fields{"module": "mq"})
	tracer = tracing.NewTracerProvider("nats").Tracer("nats")
	client = newClient()
	initMetrics()
}
