package rabbitmq

import (
  "encoding/json"

  log "github.com/sirupsen/logrus"
  amqp "github.com/rabbitmq/amqp091-go"
)

const (
  nexName = "notification-v1"
)

type NotificationChannel struct {
  *amqp.Channel
}

func InitNotificationChannel() {
  var err error

  ch, err := conn.Channel()
  if err != nil {
    log.Fatal(err)
  }

  err = ch.ExchangeDeclare(
    nexName,
    "topic",
    true,
    false,
  )
  if err != nil {
    log.Fatal(err)
  }

  channels["notification"] = NotificationChannel{ch}

}

/*
func (nc *NotificationChannel) Cancel() {
  nq.Channel.Cancel(nq.qName, false)
}
*/

type NotificationQueue struct {
  name string
  ch *NotificationChannel
}

var notificationQueues map[string]*NotificationQueue

func InitNotificationQueue(name string) *NotificationQueue {
  // Do not duplicate queues
  if nq, found := notificationQueues[name]; found {
    return nq
  }

  nq := &NotificationQueue{name, channels["notification"]}
  _, err = nq.ch.QueueDeclare(nq.Name,
    true, // durable
    false, // delete when unused
    false, // exclusive
    false, // no-wait
    nil, // args
  )
  if err != nil {
    log.Fatal(err)
  }
  err = nq.ch.QueueBind(nq.Name, "", nexName, false, nil)
  if err != nil {
    log.Fatal(err)
  }
  NotificationQueues[name] = nq
  return nq
}

type NotificationMessage struct {
  Delivery *amqp.Delivery
  Notification *v2.Notification
}

func (nq *NotificationQueue) Consume() (<-chan NotificationMessage, error) {
  out := make(chan NotificationMessage)
  dd, err := nq.ch.Consume(nq.Name, "", true, false, false, false, nil)
  if err != nil {
    return out, err
  }
  for _, d := range dd {
    var notif *v2.Notification
    if err := json.UnMarshal(d.Body, &notif); err != nil {
      log.Warn("Rejecting message (%s): %s", err, d.Body)
      d.Reject(false)
    }
    out <- NotificationMessage{d, notif}
  }
}

func (nq *NotificationQueue) Publish(ctx context.Context, notif *v2.Notification) error {
  body, err := json.Marshal(alert)
  if err != nil {
    return err
  }
  msg := amqp.Publishing{
    DeliveryMode: amqp.Persistent,
    ContextType: "application/vnd.snooze.notif.v2+json",
    Body: body,
  }
  return nq.ch.Publish(nexName, nq.Name, false, false, msg)
}
