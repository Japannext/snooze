package syslog

import (
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/health"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/rabbitmq"
)

var processChannel *rabbitmq.ProcessChannel

func Run() {

	logging.Init()
	// initConfig()
	rabbitmq.Init()
	processChannel = rabbitmq.InitProcessChannel()

	syslogServer := NewSyslogServer()
	dm := daemon.NewDaemonManager()
	dm.Add("server", syslogServer)
	h := health.HealthStatus{}
	srv := daemon.NewHttpDaemon()
	srv.Engine.GET("/livez", h.LiveRoute)
	srv.Engine.GET("/readyz", h.ReadyRoute)
	dm.Add("http", srv)

	h.Live()
	h.Ready()
	dm.Run()
}
