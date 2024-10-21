package otel

import (
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/mq"
)

var (
	processQ *mq.Pub
)

func Startup() *daemon.DaemonManager {

	logging.Init()
	initConfig()

	processQ = mq.ProcessPub()

	dm := daemon.NewDaemonManager()
	dm.AddDaemon("otel-grpc", NewOtelGrpcServer())

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
