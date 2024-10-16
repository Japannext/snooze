package silence

import (
	"github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/common/tracing"
)

var silences []*Silence
var log *logrus.Entry
var tracer = tracing.Tracer("snooze-process")

func Startup(rules []*Silence) {
	initMetrics()
	log = logrus.WithFields(logrus.Fields{"module": "silence"})
	for _, s := range rules {
		s.Load()
		silences = append(silences, s)
	}
}
