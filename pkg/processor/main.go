package processor

import (
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/utils"

	// Sub-Processors
	"github.com/japannext/snooze/pkg/processor/notification"
	"github.com/japannext/snooze/pkg/processor/profile"
	"github.com/japannext/snooze/pkg/processor/ratelimit"
	"github.com/japannext/snooze/pkg/processor/silence"
	"github.com/japannext/snooze/pkg/processor/snooze"
	"github.com/japannext/snooze/pkg/processor/store"
	"github.com/japannext/snooze/pkg/processor/transform"
)

var processQ = mq.ProcessSub()
var pool *utils.Pool

// Logic done only at the application startup.
// All errors are fatal.
func Startup() *daemon.DaemonManager {

	logging.Init()
	initConfig()
	initMetrics()
	redis.Init()

	pool = utils.NewPool(config.MaxWorkers)

	transform.Startup(pipeline.Transforms)
	silence.Startup(pipeline.Silences)
	profile.Startup(pipeline.Profiles)
	snooze.Init()
	ratelimit.Startup(pipeline.RateLimits)
	notification.Startup(pipeline.Notifications, pipeline.DefaultDestinations)
	store.Startup()

	dm := daemon.NewDaemonManager()
	dm.AddDaemon("consumer", NewConsumer())

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
