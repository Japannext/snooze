package writer

import (
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/opensearch"
)

var storeQ *mq.Sub

func Startup() *daemon.DaemonManager {

	logging.Init()
	initConfig()
	initMetrics()
	opensearch.Init()

	storeQ = mq.StoreSub()

	dm := daemon.NewDaemonManager()
	dm.AddDaemon("consumer", NewConsumer())

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
