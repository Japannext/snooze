package grouping

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var log *logrus.Entry
var tracer trace.Tracer
var groupings []*Grouping

func Startup(rules []*Grouping) {
	log = logrus.WithFields(logrus.Fields{"module": "grouping"})
	tracer = otel.Tracer("snooze")
	for _, group := range rules {
		group.Load()
		groupings = append(groupings, group)
	}
}
