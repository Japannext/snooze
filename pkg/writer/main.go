package writer

import (
	"github.com/nats-io/nats.go/jetstream"

	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/mq"
)

var storeQ jetstream.Consumer

func Startup() *daemon.DaemonManager {

	logging.Init()
	initConfig()
	initMetrics()
	mq.Init()
	storeQ = mq.Consumer(mq.STORE_STREAM)

	dm := daemon.NewDaemonManager()
	dm.AddDaemon("consumer", NewConsumer())

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
