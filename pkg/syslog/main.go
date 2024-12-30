package syslog

import (
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	processQ *mq.Pub
	tracer   trace.Tracer
)

func Startup() *daemon.DaemonManager {
	logging.Init()
	initConfig()
	initMetrics()
	mq.Init()

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
