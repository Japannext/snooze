package apiserver

import (
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/tracing"

	"github.com/japannext/snooze/pkg/apiserver/routes"

	"github.com/gin-contrib/cors"
)

// We use a route registration array so that we can
// declare routes in the same file as the route, and avoid
// having to list them all in one place.
// Routes can be listed easily enough with `git grep 'r[.].*/api/' pkg/apiserver`
// var routes []func(*gin.Engine)

func Startup() *daemon.DaemonManager {
	// Init components
	logging.Init()
	initConfig()
	routes.InitMetrics()
	opensearch.Init()
	redis.Init()
	mq.Init()

	dm := daemon.NewDaemonManager()
	srv := daemon.NewHttpDaemon(routes.Registers()...)
	srv.Engine.Use(otelgin.Middleware("snooze-alertmanager", otelgin.WithFilter(tracing.HTTPFilter)))

	// Static routes
	// srv.Engine.Group("/static", eTagMiddleware()).Static("/", config.StaticPath)

	if corsConfig != nil {
		srv.Engine.Use(cors.New(*corsConfig))
	}

	dm.AddDaemon("http", srv)

	return dm
}

func Run() {
	dm := Startup()
	dm.Run()
}
