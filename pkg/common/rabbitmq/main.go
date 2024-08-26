package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

var conn *amqp.Connection
var channelsToClose = make(map[string]ChannelInterface)

var log *logrus.Entry

func Init() {
	var err error

	log = logrus.WithFields(logrus.Fields{"module": "rabbitmq"})

	url, config := initConfig()
	Client.conn, err = amqp.DialConfig(url, config)
	if err != nil {
		log.Fatal(err)
	}
}

type ChannelInterface interface {
	Close() error
}

func Close() error {
	for _, channel := range channelsToClose {
		if err := channel.Close(); err != nil {
			return err
		}
	}
	return conn.Close()
}
