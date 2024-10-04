package alertmanager

import (
    "github.com/japannext/snooze/pkg/common/daemon"
    "github.com/japannext/snooze/pkg/common/logging"
    "github.com/japannext/snooze/pkg/common/redis"
    "github.com/japannext/snooze/pkg/common/opensearch"
)

func Startup() *daemon.DaemonManager {

    logging.Init()
    initConfig()
    initMetrics()
	redis.Init()
	opensearch.Init()

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
