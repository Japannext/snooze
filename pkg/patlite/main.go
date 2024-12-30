package patlite

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/tracing"
)

var (
	notifyQ *mq.Sub
	storeQ  *mq.Pub
	tracer  trace.Tracer
)

func Startup() *daemon.DaemonManager {
	logging.Init()
	initConfig()
	loadProfiles()
	tracing.Init("snooze-googlechat")
	mq.Init()

	// client.TestSpaces()

	tracer = otel.Tracer("snooze")

	notifyQ = mq.NotifySub("googlechat")
	storeQ = mq.StorePub()

	dm := daemon.NewDaemonManager()

	wp := mq.NewWorkerPool(notifyQ, notificationHandler, 20)
	dm.AddDaemon("notifier", wp)

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
