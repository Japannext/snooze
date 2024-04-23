package rabbitmq

import (
  log "github.com/sirupsen/logrus"
  amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var channelsToClose map[string]ChannelInterface

func Init() {
  var err error

  config := initConfig()

  conn, err = amqp.Dial(config.Address)
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
