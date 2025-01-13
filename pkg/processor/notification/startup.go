package notification

import (
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/models"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	notifications       []*Notification
	defaultDestinations []models.Destination
	log                 *logrus.Entry
	tracer              trace.Tracer
	notifyQ, storeQ     *mq.Pub
)

func Startup(notifs []*Notification, defaults []models.Destination) {
	log = logrus.WithFields(logrus.Fields{"module": "notification"})
	tracer = otel.Tracer("snooze")

	notifyQ = mq.NotifyPub()
	storeQ = mq.StorePub().WithIndex(models.NotificationIndex)

	defaultDestinations = defaults

	for _, notif := range notifs {
		notif.Load()
		notifications = append(notifications, notif)
	}
}
