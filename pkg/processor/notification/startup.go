package notification

import (
	"github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/rabbitmq"
	"github.com/japannext/snooze/pkg/common/utils"
)

var notifications []*Notification
var defaultDestinations []models.Destination
var log *logrus.Entry
var producers = map[string]*rabbitmq.Producer{}

func Startup(notifs []*Notification, defaults []models.Destination) {
	log = logrus.WithFields(logrus.Fields{"module": "notification"})

	defaultDestinations = defaults

	var queues = utils.NewOrderedSet[string]()

	for _, notif := range notifs {
		notif.Load()
		notifications = append(notifications, notif)
		for _, dest := range notif.Destinations {
			if dest.Queue == "dummy" {
				continue
			}
			queues.Append(dest.Queue)
		}
	}

	for _, queue := range queues.Items() {
		producers[queue] = rabbitmq.NewNotificationProducer(queue)
	}
}
