package processor

import (
	"github.com/nats-io/nats.go/jetstream"

	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/redis"

	// Sub-Processors
	"github.com/japannext/snooze/pkg/processor/notification"
	"github.com/japannext/snooze/pkg/processor/profile"
	"github.com/japannext/snooze/pkg/processor/ratelimit"
	"github.com/japannext/snooze/pkg/processor/silence"
	"github.com/japannext/snooze/pkg/processor/snooze"
	"github.com/japannext/snooze/pkg/processor/store"
	"github.com/japannext/snooze/pkg/processor/transform"
	"github.com/japannext/snooze/pkg/processor/tracing"
)

var processQ jetstream.Consumer

// Logic done only at the application startup.
// All errors are fatal.
func Startup() *daemon.DaemonManager {

	logging.Init()
	initConfig()
	tracing.Init()
	initMetrics()
	redis.Init()
	mq.Init()

	transform.Startup(pipeline.Transforms)
	silence.Startup(pipeline.Silences)
	profile.Startup(pipeline.Profiles)
	snooze.Init()
	ratelimit.Startup(pipeline.RateLimits)
	notification.Startup(pipeline.Notifications, pipeline.DefaultDestinations)
	store.Startup()

	processQ = mq.Consumer(mq.PROCESS_STREAM)

	dm := daemon.NewDaemonManager()
	dm.AddDaemon("consumer", NewConsumer())

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
