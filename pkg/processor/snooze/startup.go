package snooze

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

func Init() {
	tracer = otel.Tracer("snooze")
	log = logrus.WithFields(logrus.Fields{"module": "snooze"})
	storeQ = mq.StorePub()
	initMetrics()
}
