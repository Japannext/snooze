package transform

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	log        *logrus.Entry
	tracer     trace.Tracer
	transforms []*Transform
)

func Startup(trs []*Transform) {
	log = logrus.WithFields(logrus.Fields{"module": "transform"})
	tracer = otel.Tracer("snooze")
	for _, tr := range trs {
		tr.Load()
		transforms = append(transforms, tr)
	}
}
