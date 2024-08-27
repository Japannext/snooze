package apiserver

import (
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/redis"

	"github.com/gin-contrib/cors"
)

func Startup() *daemon.DaemonManager {
	// Init components
	logging.Init()
	initConfig()
	opensearch.Init()
	redis.Init()

	dm := daemon.NewDaemonManager()
	srv := daemon.NewHttpDaemon()

	// Static routes
	// srv.Engine.Group("/static", eTagMiddleware()).Static("/", config.StaticPath)

	if corsConfig != nil {
	srv.Engine.Use(cors.New(*corsConfig))
	}

	// Dynamic routes
	registerAuthRoutes(srv.Engine)
	registerLogRoutes(srv.Engine)
	registerNotificationRoutes(srv.Engine)
	registerSampleRoute(srv.Engine)
	dm.AddDaemon("http", srv)

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
