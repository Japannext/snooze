package rabbitmq

const (
	NOTIFICATION_EXCHANGE = "notification-v2"
	NOTIFICATION_CONTENT_TYPE = "application/vnd.snooze.notification.v2+json"
)

func SetupNotifications(queueNames []string) {
	// Setup
	exchanges := map[string]ExchangeOptions{
		NOTIFICATION_EXCHANGE: {Kind: "topic", Durable: true},
	}
	queues := map[string]QueueOptions{}
	binds := map[string]BindOptions{}
	for _, name := range queueNames {
		queues[name] = QueueOptions{Durable: true}
		binds[name] = BindOptions{
			Exchange: NOTIFICATION_EXCHANGE,
			Key: name,
		}
	}
	Client.setup(exchanges, queues, binds)
}

func NewNotificationProducer(topic string) *Producer {
	return NewProducer(NOTIFICATION_EXCHANGE, topic, NOTIFICATION_CONTENT_TYPE)
}
