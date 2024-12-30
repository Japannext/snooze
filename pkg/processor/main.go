package processor

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/tracing"
	"github.com/japannext/snooze/pkg/common/utils"
	// Sub-Processors
	"github.com/japannext/snooze/pkg/processor/grouping"
	"github.com/japannext/snooze/pkg/processor/mapping"
	"github.com/japannext/snooze/pkg/processor/notification"
	"github.com/japannext/snooze/pkg/processor/profile"
	"github.com/japannext/snooze/pkg/processor/ratelimit"
	"github.com/japannext/snooze/pkg/processor/silence"
	"github.com/japannext/snooze/pkg/processor/snooze"
	"github.com/japannext/snooze/pkg/processor/store"
	"github.com/japannext/snooze/pkg/processor/transform"
)

var processQ *mq.Sub
var pool *utils.Pool
var tracer trace.Tracer

var routes []func(*gin.Engine)

// Logic done only at the application startup.
// All errors are fatal.
func Startup() *daemon.DaemonManager {

	logging.Init()
	initConfig()
	initMetrics()
	redis.Init()
	tracing.Init("snooze-process")
	mq.Init()

	processQ = mq.ProcessSub()
	pool = utils.NewPool(config.MaxWorkers)
	tracer = otel.Tracer("snooze")

	mapping.Startup(pipeline.Mappings)
	transform.Startup(pipeline.Transforms)
	grouping.Startup(pipeline.Grouping)
	silence.Startup(pipeline.Silences)
	profile.Startup(pipeline.Profiles)
	snooze.Init()
	ratelimit.Startup(pipeline.RateLimits)
	notification.Startup(pipeline.Notifications, pipeline.DefaultDestinations)
	store.Startup()

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
