package rabbitmq

import (
  log "github.com/sirupsen/logrus"
  amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
  name string
  conn *amqp.Connection
  ch *amqp.Channel
  *amqp.Queue
}

func (q *Queue) Init() {
  conn, err := amqp.Dial("")
  if err != nil {
    log.Fatal(err)
  }
  defer conn.Close()
}

func (q *Queue) Close() {
  q.channel.Close()
  q.connection.Close()
}

func initq() *Queue {
  conn, err := amqp.Dial("...")
  if err != nil {
    log.Fatal(err)
  }
  ch, err := conn.Channel()
  if err != nil {
    log.Fatal(err)
  }
  return &Queue{conn: conn, ch: ch}
}

func ProcessQueue() *Queue {
  q := initq()
  qq, err := q.ch.QueueDeclare("processing", false, false, false, false, nil)
  q.q = qq
}

q, err := ch.QueueDeclare(
  "processing",
  false, // durable
  false, // delete when unused
  false, // exclusive
  false, // no-wait
  nil, // args
)
