package rabbitmq

import (
  amqp "github.com/rabbitmq/amqp091-go"
)

type exOpts struct {
  Durable bool
  AutoDelete bool
  Internal bool
  NoWait bool
  Args map[string]interface{}
}

func exchange(ch *amqp.Channel, name, kind string, opts *exOpts) error {
  return ch.ExchangeDeclare(name, kind, opts.Durable, opts.AutoDelete, opts.Internal, opts.NoWait, opts.Args)
}

type qOpts struct {
  Durable bool
  DeleteWhenUnused bool
  Exclusive bool
  NoWait bool
  Args map[string]interface{}
}

func queue(ch *amqp.Channel, name string, opts *qOpts) (amqp.Queue, error) {
  return ch.QueueDeclare(name, opts.Durable, opts.DeleteWhenUnused, opts.Exclusive, opts.NoWait, opts.Args)
}

type chOpts struct {
  AutoAck bool
  Exclusive bool
  NoLocal bool
  NoWait bool
  Args map[string]interface{}
}

func consume(ch *amqp.Channel, q, topic string, opts *chOpts) (<-chan amqp.Delivery, error) {
  return ch.Consume(q, topic, opts.AutoAck, opts.Exclusive, opts.NoLocal, opts.NoWait, opts.Args)
}
