package alertcloser

import (
	"time"

	"github.com/japannext/snooze/internal/daemon"
	"github.com/japannext/snooze/internal/logging"
	"github.com/japannext/snooze/internal/mq"
	"github.com/japannext/snooze/internal/opensearch"
	"github.com/japannext/snooze/internal/redis"
	"github.com/japannext/snooze/internal/tracing"
)

var (
	storeQ *mq.Pub
	//tracer trace.Tracer
)

const closeExpiredAlertsInterval = 10 * time.Second

func Startup() *daemon.DaemonManager {
	logging.Init()
	// initConfig()
	// initMetrics()
	redis.Init()
	opensearch.Init()
	mq.Init()

	tracing.Init("snooze-alertmanager")
	// tracer = otel.Tracer("snooze")

	storeQ = mq.StorePub()

	dm := daemon.NewDaemonManager()

	closeExpiredAlertsDaemon := daemon.NewLockDaemon("job:close-expired-alerts", closeExpiredAlertsInterval, closeExpiredAlerts)
	dm.AddDaemon("close-expired-alerts", closeExpiredAlertsDaemon)

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
