package alertmanager

import (
	"time"

	"github.com/japannext/snooze/internal/daemon"
	"github.com/japannext/snooze/internal/logging"
	"github.com/japannext/snooze/internal/mq"
	"github.com/japannext/snooze/internal/opensearch"
	"github.com/japannext/snooze/internal/redis"
	"github.com/japannext/snooze/internal/tracing"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	storeQ *mq.Pub
	tracer trace.Tracer
)

const closeExpiredAlertsInterval = 10 * time.Second

func Startup() *daemon.DaemonManager {
	logging.Init()
	initConfig()
	initMetrics()
	redis.Init()
	opensearch.Init()
	mq.Init()

	tracing.Init("snooze-alertmanager")
	tracer = otel.Tracer("snooze")

	storeQ = mq.StorePub()

	dm := daemon.NewDaemonManager()
	srv := daemon.NewHttpDaemon()
	srv.Engine.Use(otelgin.Middleware("snooze-alertmanager", otelgin.WithFilter(tracing.HTTPFilter)))
	{
		srv.Engine.POST("/api/v2/alerts", postAlert)
	}
	dm.AddDaemon("http", srv)

	closeExpiredAlertsDaemon := daemon.NewLockDaemon("job:close-expired-alerts", closeExpiredAlertsInterval, closeExpiredAlerts)
	dm.AddDaemon("close-expired-alerts", closeExpiredAlertsDaemon)

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
