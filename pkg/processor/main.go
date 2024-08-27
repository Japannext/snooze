package processor

import (
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/rabbitmq"
	"github.com/japannext/snooze/pkg/common/redis"

	// Sub-Processors
	"github.com/japannext/snooze/pkg/processor/profile"
	"github.com/japannext/snooze/pkg/processor/notification"
	"github.com/japannext/snooze/pkg/processor/ratelimit"
	"github.com/japannext/snooze/pkg/processor/silence"
	"github.com/japannext/snooze/pkg/processor/snooze"
	"github.com/japannext/snooze/pkg/processor/transform"
)

// Logic done only at the application startup.
// All errors are fatal.
func Startup() *daemon.DaemonManager {

	logging.Init()
	initConfig()
	opensearch.Init()
	redis.Init()
	rabbitmq.Init()

	transform.InitRules(pipeline.TransformRules)
	silence.InitRules(pipeline.SilenceRules)
	profile.InitRules(pipeline.Profiles)
	snooze.Init()
	ratelimit.Init(pipeline.RateLimit)
	notification.Startup(pipeline.NotificationRules, pipeline.DefaultNotificationChannels)

	dm := daemon.NewDaemonManager()
	dm.AddDaemon("processor", NewProcessor())

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
