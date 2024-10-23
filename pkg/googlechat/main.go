package googlechat

import (
    "github.com/japannext/snooze/pkg/common/daemon"
    "github.com/japannext/snooze/pkg/common/logging"
    "github.com/japannext/snooze/pkg/common/mq"
    "github.com/japannext/snooze/pkg/models"
)

var (
	notifyQ *mq.Sub
	storeQ *mq.Pub
)

func Startup() *daemon.DaemonManager {
    logging.Init()
    initConfig()
    loadProfiles()
	initGooglechat()
	mq.Init()

	notifyQ = mq.NotifySub("googlechat")
	storeQ = mq.StorePub().WithIndex(models.NOTIFICATION_INDEX)

    dm := daemon.NewDaemonManager()

	wp := mq.NewWorkerPool(notifyQ, handler, 20)
    dm.AddDaemon("notifier", wp)

    return dm
}

func Run() {
    dm := Startup()
    dm.Run()
}
