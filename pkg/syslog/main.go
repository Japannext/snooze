package syslog

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/tracing"
)

var processQ *mq.Pub
var tracer trace.Tracer

func Startup() *daemon.DaemonManager {

	logging.Init()
	initConfig()
	initMetrics()

	processQ = mq.ProcessPub()
	tracing.Init("snooze-syslog")
	tracer = otel.Tracer("snooze")

	dm := daemon.NewDaemonManager()
	dm.AddDaemon("syslog", NewSyslogServer())

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
