package store

import (
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	log    *logrus.Entry
	tracer trace.Tracer
	storeQ *mq.Pub
)

func Startup() {
	log = logrus.WithFields(logrus.Fields{"module": "store"})
	tracer = otel.Tracer("snooze")
	storeQ = mq.StorePub()

	initMetrics()
}
