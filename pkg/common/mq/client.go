package mq

import (
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"module": "mq"})

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

var clientInstance *Client
var clientOnce sync.Once

func getClient() *Client {
	clientOnce.Do(func() {
		clientInstance = newClient()
	})
	return clientInstance
}
