package apiserver

import (
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/health"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/redis"

	"github.com/gin-contrib/cors"
)

func Run() {
	// Init components
	logging.Init()
	health.Init()
	initConfig()
	opensearch.Init()
	redis.Init()

	dm := daemon.NewDaemonManager()
	srv := daemon.NewHttpDaemon()

	// Health
	h := health.HealthStatus{}
	srv.Engine.GET("/livez", h.LiveRoute)
	srv.Engine.GET("/readyz", h.ReadyRoute)

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
	dm.Add("http", srv)

	h.Live()
	h.Ready()

	// Listen
	// hostport := fmt.Sprintf("%s:%d", config.ListenAddress, config.ListenPort)
	dm.Run()
}
