package activecheck

import (
	"github.com/japannext/snooze/internal/daemon"
	"github.com/japannext/snooze/internal/logging"
)

func Startup() *daemon.DaemonManager {
	logging.Init()
	initConfig()
	// initSyslogMetrics()
	dm := daemon.NewDaemonManager()
	srv := daemon.NewHttpDaemon()
	srv.Engine.POST("/webhook/:uid", callbackHandler)
	srv.Engine.GET("/probe/:name", probeHandler)
	dm.AddDaemon("http", srv)
	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
