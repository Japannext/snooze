package mail

import (
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/health"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/rabbitmq"
	"github.com/japannext/snooze/pkg/common/notifier"
)

func Startup() *daemon.DaemonManager {
	logging.Init()
	initConfig()
	loadProfiles()
	rabbitmq.Init()

	dm := daemon.NewDaemonManager()

	notifier := notifier.NewNotifier(config.Queue, handler)
	dm.Add("notifier", notifier)

	h := health.HealthStatus{}
	dm.Add("health", health.NewHealthDaemon(h))

	h.Live()
	h.Ready()
	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
