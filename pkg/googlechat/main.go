package googlechat

import (
    "github.com/japannext/snooze/pkg/common/daemon"
    "github.com/japannext/snooze/pkg/common/logging"
    "github.com/japannext/snooze/pkg/common/rabbitmq"
    "github.com/japannext/snooze/pkg/common/notifier"
)

func Startup() *daemon.DaemonManager {
    logging.Init()
    initConfig()
    loadProfiles()
    rabbitmq.Init()
	initGooglechat()

    dm := daemon.NewDaemonManager()

    notifier := notifier.NewNotifier("googlechat", handler)
    dm.AddDaemon("notifier", notifier)
    dm.AddReadyCheck(notifier.Consumer)

    return dm
}

func Run() {
    dm := Startup()
    dm.Run()
}
