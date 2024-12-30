package alertmanager

import (
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/tracing"
	"github.com/japannext/snooze/pkg/models"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	storeQ *mq.Pub
	tracer trace.Tracer
)

func Startup() *daemon.DaemonManager {
	logging.Init()
	initConfig()
	initMetrics()
	redis.Init()
	mq.Init()

	tracing.Init("snooze-alertmanager")
	tracer = otel.Tracer("snooze")

	storeQ = mq.StorePub().WithIndex(models.ALERT_INDEX)

	dm := daemon.NewDaemonManager()
	srv := daemon.NewHttpDaemon()
	srv.Engine.Use(otelgin.Middleware("snooze-alertmanager", otelgin.WithFilter(tracing.HTTPFilter)))
	{
		srv.Engine.POST("/api/v2/alerts", postAlert)
	}
	dm.AddDaemon("http", srv)

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
