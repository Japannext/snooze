package processor

import (
	"github.com/gin-gonic/gin"
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/tracing"
	"github.com/japannext/snooze/pkg/common/utils"
	"github.com/japannext/snooze/pkg/processor/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	processQ *mq.Sub
	pool     *utils.Pool
	tracer   trace.Tracer
)

var routes []func(*gin.Engine)

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
	srv := daemon.NewHttpDaemon(routes...)
	dm.AddDaemon("http", srv)
	dm.AddDaemon("consumer", NewConsumer())

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
