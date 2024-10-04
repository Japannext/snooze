package syslog

import (
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/rabbitmq"
)

var producer *rabbitmq.Producer

func Startup() *daemon.DaemonManager {

	logging.Init()
	initConfig()
	initMetrics()
	rabbitmq.Init()
	producer = rabbitmq.NewLogProducer()

	dm := daemon.NewDaemonManager()
	dm.AddDaemon("server", NewSyslogServer())

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
