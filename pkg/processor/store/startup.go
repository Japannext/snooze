package store

import (
	"github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/common/tracing"
)

var log *logrus.Entry
var tracer = tracing.Tracer("snooze-process")

func Startup() {
	log = logrus.WithFields(logrus.Fields{"module": "store"})
	initMetrics()
}
