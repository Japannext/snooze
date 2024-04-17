package rabbitmq

import (
  log "github.com/sirupsen/logrus"
  amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var channels map[string]ChannelInterface

func Init() {
  var err error

  cfg := initConfig()

  conn, err = amqp.Dial(config.Address)
  if err != nil {
    log.Fatal(err)
  }
}

type ChannelInterface interface {
  Close()
}

func Close() {
  for _, channel := range channels {
    channel.Close()
  }
  conn.Close()
}
