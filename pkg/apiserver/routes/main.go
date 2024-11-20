package routes

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/tracing"
)

var storeQ *mq.Pub
var tracer trace.Tracer

var registers []func(*gin.Engine)

func Registers() []func(*gin.Engine) {
	storeQ = mq.StorePub()
	tracing.Init("snooze-apiserver")
	tracer = otel.Tracer("snooze")
	return registers
}
