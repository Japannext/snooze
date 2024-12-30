package silence

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	silences []*Silence
	log      *logrus.Entry
	tracer   trace.Tracer
)

func Startup(rules []*Silence) {
	initMetrics()
	log = logrus.WithFields(logrus.Fields{"module": "silence"})
	tracer = otel.Tracer("snooze")
	for _, s := range rules {
		s.Load()
		silences = append(silences, s)
	}
}
