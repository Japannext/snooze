package writer

import (
	"github.com/japannext/snooze/internal/daemon"
	"github.com/japannext/snooze/internal/logging"
	"github.com/japannext/snooze/internal/mq"
	"github.com/japannext/snooze/internal/opensearch"
	"github.com/japannext/snooze/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	storeQ *mq.Sub
	tracer trace.Tracer
)

func Startup() *daemon.DaemonManager {
	logging.Init()
	initConfig()
	initMetrics()
	opensearch.Init()
	tracing.Init("snooze-writer")
	mq.Init()

	storeQ = mq.StoreSub()
	tracer = otel.Tracer("snooze")

	dm := daemon.NewDaemonManager()
	dm.AddDaemon("consumer", NewConsumer())

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
