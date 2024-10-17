package snooze

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var log logrus.Entry
var tracer trace.Tracer

func Init() {
	tracer = otel.Tracer("snooze")
}
