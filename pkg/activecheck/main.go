package activecheck

import (
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
)

func Startup() *daemon.DaemonManager {
	logging.Init()
	initConfig()
	// initSyslogMetrics()
	dm := daemon.NewDaemonManager()
	srv := daemon.NewHttpDaemon()
	srv.Engine.POST("/webhook/:uid", webhookHandler)
	srv.Engine.GET("/probe/:name", probeHandler)
	dm.AddDaemon("http", srv)
	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
