package mail

import (
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/notifier"
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

	notifyQ = mq.NotifySub("mail")
	storeQ = mq.StorePub().WithIndex(models.NOTIFICATION_INDEX)

	dm := daemon.NewDaemonManager()

	notifier := notifier.NewNotifier(config.Queue, handler)
	dm.AddDaemon("notifier", notifier)
	dm.AddReadyCheck(notifier.Consumer)

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
