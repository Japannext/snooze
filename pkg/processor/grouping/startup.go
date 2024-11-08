package grouping

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/japannext/snooze/pkg/common/mq"
)

var log *logrus.Entry
var tracer trace.Tracer
var groupings []*Grouping
var storeQ *mq.Pub

func Startup(rules []*Grouping) {
	log = logrus.WithFields(logrus.Fields{"module": "grouping"})
	tracer = otel.Tracer("snooze")
	storeQ = mq.StorePub()
	for _, group := range rules {
		group.Load()
		groupings = append(groupings, group)
	}
}
