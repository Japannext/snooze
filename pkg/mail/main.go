package mail

import (
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/models"
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
	storeQ = mq.StorePub().WithIndex(models.NOTIFICATION_INDEX)

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
