package processor

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/japannext/snooze/components/processor/metrics"
	"github.com/japannext/snooze/internal/daemon"
	"github.com/japannext/snooze/internal/logging"
	"github.com/japannext/snooze/internal/mq"
	"github.com/japannext/snooze/internal/redis"
	"github.com/japannext/snooze/internal/tracing"
	"github.com/japannext/snooze/internal/utils"
)

var (
	processQ *mq.Sub
	pool     *utils.Pool
	tracer   trace.Tracer
)

// Logic done only at the application startup.
// All errors are fatal.
func Startup() *daemon.DaemonManager {
	logging.Init()
	initConfig()
	metrics.Init()
	redis.Init()
	tracing.Init("snooze-process")
	mq.Init()

	processQ = mq.ProcessSub()
	pool = utils.NewPool(env.MaxWorkers)
	tracer = otel.Tracer("snooze")

	dm := daemon.NewDaemonManager()
	srv := daemon.NewHttpDaemon()
	dm.AddDaemon("http", srv)
	dm.AddDaemon("consumer", NewConsumer())

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
