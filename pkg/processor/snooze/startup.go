package snooze

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/japannext/snooze/pkg/common/mq"
)

var log *logrus.Entry
var tracer trace.Tracer
var storeQ *mq.Pub

func Init() {
	tracer = otel.Tracer("snooze")
	log = logrus.WithFields(logrus.Fields{"module": "snooze"})
	storeQ = mq.StorePub()
	initMetrics()
}
