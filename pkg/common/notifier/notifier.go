package notifier

import (
	"encoding/json"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/rabbitmq"
)

const (
	EXCHANGE_NAME = "notification-v2"
)

type NotificationHandler = func(*api.Notification) error

type Notifier struct {
	Queue string
	Consumer *rabbitmq.Consumer
	handler NotificationHandler
}

func NewNotifier(queueName string, handler NotificationHandler) *Notifier {
	notifier := &Notifier{Queue: queueName, handler: handler}

	rabbitmq.SetupNotifications([]string{queueName})

	options := rabbitmq.ConsumerOptions{AutoAck: true}
	notifier.Consumer = rabbitmq.NewConsumer(queueName, queueName, transform(handler), options)

	return notifier
}

func transform(handler NotificationHandler) rabbitmq.Handler {
	return func(delivery rabbitmq.Delivery) error {
		var notification *api.Notification
		if err := json.Unmarshal(delivery.Body, &notification); err != nil {
			return err
		}
		if err := handler(notification); err != nil {
			return err
		}
		return nil
	}
}

func (n *Notifier) Run() error {
	return n.Consumer.ConsumeForever()
}

func (n *Notifier) Stop() {
	n.Consumer.Stop()
}
