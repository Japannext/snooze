package apiserver

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/tracing"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

var tracer trace.Tracer

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

	tracing.Init("snooze-apiserver")
	tracer = otel.Tracer("snooze")

	dm := daemon.NewDaemonManager()
	srv := daemon.NewHttpDaemon()
	srv.Engine.Use(otelgin.Middleware("snooze-alertmanager", otelgin.WithFilter(tracing.HTTPFilter)))

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
