package mail

import (
	"github.com/japannext/snooze/internal/daemon"
	"github.com/japannext/snooze/internal/logging"
	"github.com/japannext/snooze/internal/mq"
	"github.com/japannext/snooze/internal/models"
)

var (
	notifyQ *mq.Sub
	storeQ  *mq.Pub
)

func Startup() *daemon.DaemonManager {
	logging.Init()
	initConfig()
	loadProfiles()
	mq.Init()

	notifyQ = mq.NotifySub("mail")
	storeQ = mq.StorePub().WithIndex(models.NotificationIndex)

	dm := daemon.NewDaemonManager()

	// Worker pool
	wp := mq.NewWorkerPool(notifyQ, handler, 20)
	dm.AddDaemon("worker_pool", wp)

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
