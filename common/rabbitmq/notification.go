package rabbitmq

import (
  "context"
  "encoding/json"

  log "github.com/sirupsen/logrus"
  amqp "github.com/rabbitmq/amqp091-go"

  api "github.com/japannext/snooze/common/api/v2"
)

const (
  nexName = "notification-v1"
)

type NotificationChannel struct {
  *amqp.Channel
}

func InitNotificationChannel() *NotificationChannel {
  var err error

  ch, err := conn.Channel()
  if err != nil {
    log.Fatal(err)
  }

  err = exchange(ch, nexName, "topic", &exOpts{Durable: true})
  if err != nil {
    log.Fatal(err)
  }

  not := &NotificationChannel{ch}
  channelsToClose["notification"] = not
  return not
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

/*
Example usage for consumer:
var nc *rabbitmq.NotificationChannel
var nq *rabbitmq.NotificationQueue
func init() {
  nc = InitNotificationChannel()
  nq = nc.NewQueue("myplugin")
}
func run() {
  notifs, err := nc.Consume()
  if err != nil {
    log.Error(err)
  }
  for notif := range notif {
    // handle notification...
  }
}
*/
func (nc *NotificationChannel) NewQueue(name string) *NotificationQueue {
  // Do not duplicate queues
  if nq, found := notificationQueues[name]; found {
    return nq
  }

  ch := nc.Channel
  nq := &NotificationQueue{name, nc}
  _, err := queue(ch, nq.name, &qOpts{Durable: true})
  if err != nil {
    log.Fatal(err)
  }
  err = ch.QueueBind(nq.name, "", nexName, false, nil)
  if err != nil {
    log.Fatal(err)
  }
  return nq
}

type NotificationMessage struct {
  Delivery *amqp.Delivery
  Notification *api.Notification
}

func (nq *NotificationQueue) Consume() (<-chan NotificationMessage, error) {
  out := make(chan NotificationMessage)
  dd, err := nq.ch.Consume(nq.name, "", true, false, false, false, nil)
  if err != nil {
    return out, err
  }
  for d := range dd {
    var notif *api.Notification
    if err := json.Unmarshal(d.Body, &notif); err != nil {
      log.Warnf("Rejecting message (%s): %s", err, d.Body)
      d.Reject(false)
    }
    out <- NotificationMessage{&d, notif}
  }
  return out, nil
}

func (nq *NotificationQueue) Publish(ctx context.Context, notif *api.Notification) error {
  body, err := json.Marshal(notif)
  if err != nil {
    return err
  }
  msg := amqp.Publishing{
    DeliveryMode: amqp.Persistent,
    ContentType: "application/vnd.snooze.notif.v2+json",
    Body: body,
  }
  return nq.ch.Publish(nexName, nq.name, false, false, msg)
}
