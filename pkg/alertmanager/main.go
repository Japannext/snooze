package alertmanager

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

    "github.com/japannext/snooze/pkg/common/daemon"
    "github.com/japannext/snooze/pkg/common/logging"
    "github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/tracing"
)

var storeQ *mq.Pub
var tracer trace.Tracer

func Startup() *daemon.DaemonManager {

    logging.Init()
    initConfig()
    initMetrics()
	redis.Init()

	tracing.Init("snooze-alertmanager")
	tracer = otel.Tracer("snooze")

	storeQ = mq.StorePub().WithIndex(models.ALERT_INDEX)

    dm := daemon.NewDaemonManager()
	srv := daemon.NewHttpDaemon()
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
