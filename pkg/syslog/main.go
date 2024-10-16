package syslog

import (
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/tracing"
)

var processQ = mq.ProcessPub()
var tracer = tracing.Tracer("snooze-syslog")

func Startup() *daemon.DaemonManager {

	logging.Init()
	initConfig()
	initMetrics()

	dm := daemon.NewDaemonManager()
	dm.AddDaemon("syslog", NewSyslogServer())

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
