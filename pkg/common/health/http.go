package health

import (
	"github.com/japannext/snooze/pkg/common/daemon"
)

func NewHealthDaemon(h HealthStatus) *daemon.HttpDaemon {
	srv := daemon.NewHttpDaemon()
	srv.Engine.GET("/livez", h.LiveRoute)
	srv.Engine.GET("/readyz", h.ReadyRoute)
	return srv
}
