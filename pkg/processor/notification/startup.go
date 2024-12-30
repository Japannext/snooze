package notification

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/models"
)

var notifications []*Notification
var defaultDestinations []models.Destination
var log *logrus.Entry
var tracer trace.Tracer
var notifyQ, storeQ *mq.Pub

func Startup(notifs []*Notification, defaults []models.Destination) {
	log = logrus.WithFields(logrus.Fields{"module": "notification"})
	tracer = otel.Tracer("snooze")

	notifyQ = mq.NotifyPub()
	storeQ = mq.StorePub().WithIndex(models.NOTIFICATION_INDEX)

	defaultDestinations = defaults

	for _, notif := range notifs {
		notif.Load()
		notifications = append(notifications, notif)
	}
}
