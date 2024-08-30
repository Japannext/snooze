package apiserver

import (
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/redis"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

// We use a route registration array so that we can
// declare routes in the same file as the route, and avoid
// having to list them all in one place.
// Routes can be listed easily enough with `git grep 'r[.].*/api/' pkg/apiserver`
type RouteRegistration = func(*gin.Engine)
var routes []RouteRegistration

func Startup() *daemon.DaemonManager {
	// Init components
	logging.Init()
	initConfig()
	initMetrics()
	opensearch.Init()
	redis.Init()

	dm := daemon.NewDaemonManager()
	srv := daemon.NewHttpDaemon()

	// Static routes
	// srv.Engine.Group("/static", eTagMiddleware()).Static("/", config.StaticPath)

	if corsConfig != nil {
	srv.Engine.Use(cors.New(*corsConfig))
	}

	// Registering routes
	for _, register := range routes {
		register(srv.Engine)
	}
	dm.AddDaemon("http", srv)

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
