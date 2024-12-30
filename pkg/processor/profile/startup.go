package profile

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	fastMapper *FastMapper
	log        *logrus.Entry
	tracer     trace.Tracer
)

func Startup(prfs []*Profile) {
	log = logrus.WithFields(logrus.Fields{"module": "profile"})
	tracer = otel.Tracer("snooze")
	fastMapper = NewFastMapper(prfs)
}
