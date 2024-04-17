package rabbitmq

import (
  "encoding/json"

  log "github.com/sirupsen/logrus"
  amqp "github.com/rabbitmq/amqp091-go"
)

const (
  pexName = "processing-v1"
  pqName = "processing-v1"
)

type ProcessChannel struct {
  *amqp.Channel
}

func InitProcessChannel() {
  ch, err := conn.Channel()
  if err != nil {
    log.Fatal(err)
  }
  err = ch.ExchangeDeclare(nq.exName,
    "direct",
    true, // durable
    false, // auto-deleted
    false, // internal
    false, // no-wait
    nil, // args
  )
  if err != nil {
    log.Fatal(err)
  }

  _, err = ch.QueueDeclare(nq.qName,
    true, // durable
    false, // delete when unused
    false, // exclusive
    false, // no-wait
    nil, // args
  )
  if err != nil {
    log.Fatal(err)
  }
  err = ch.QueueBind(nq.qName, "", nq.exName, false, nil)
  if err != nil {
    log.Fatal(err)
  }

  channels["process"] = &ProcessChannel{ch}
}

func (pc *ProcessChannel) Cancel() {
  pc.Cancel(pqName, false)
}

type AlertMessage struct {
  Delivery *amqp.Delivery
  Alert *v2.Alert
}

func (nq *ProcessChannel) Consume() (<-chan AlertMessage, error) {
  out := make(chan AlertMessage)
  dd, err := pc.Consume(pqName, "", true, false, false, false, nil)
  if err != nil {
    return out, err
  }
  for _, d := range dd {
    var alert *v2.Alert
    if err := json.UnMarshal(d.Body, &alert); err != nil {
      log.Warn("Rejecting message (%s): %s", err, d.Body)
      d.Reject(false)
    }
    out <- AlertMessage{d, alert}
  }
}

func (pc *ProcessChannel) Publish(ctx context.Context, alert *v2.Alert) error {
  body, err := json.Marshal(alert)
  if err != nil {
    return err
  }
  msg := amqp.Publishing{
    DeliveryMode: amqp.Persistent,
    ContextType: "application/vnd.snooze.alert.v2+json",
    Body: body,
  }
  return pc.Publish(pqName, "", false, false, msg)
}
