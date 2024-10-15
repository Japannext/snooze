package mq

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
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

func Init() {
	SetupClient()
	SetupStreams()
}
