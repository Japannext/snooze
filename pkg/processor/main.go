package processor

import (
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/rabbitmq"
	"github.com/japannext/snooze/pkg/common/redis"

	// Sub-Processors
	"github.com/japannext/snooze/pkg/processor/notification"
	"github.com/japannext/snooze/pkg/processor/profile"
	"github.com/japannext/snooze/pkg/processor/ratelimit"
	"github.com/japannext/snooze/pkg/processor/silence"
	"github.com/japannext/snooze/pkg/processor/snooze"
	"github.com/japannext/snooze/pkg/processor/store"
	"github.com/japannext/snooze/pkg/processor/transform"
)

// Logic done only at the application startup.
// All errors are fatal.
func Startup() *daemon.DaemonManager {

	logging.Init()
	initConfig()
	initMetrics()
	opensearch.Init()
	redis.Init()
	rabbitmq.Init()

	transform.Startup(pipeline.Transforms)
	silence.Startup(pipeline.Silences)
	profile.Startup(pipeline.Profiles)
	snooze.Init()
	ratelimit.Startup(pipeline.RateLimits)
	notification.Startup(pipeline.Notifications, pipeline.DefaultDestinations)
	store.Startup()

	dm := daemon.NewDaemonManager()
	dm.AddDaemon("processor", NewProcessor())

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
