package store

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/japannext/snooze/pkg/common/mq"
)

var log *logrus.Entry
var tracer trace.Tracer
var storeQ *mq.Pub

func Startup() {
	log = logrus.WithFields(logrus.Fields{"module": "store"})
	tracer = otel.Tracer("snooze")
	storeQ = mq.StorePub()

	initMetrics()
}
