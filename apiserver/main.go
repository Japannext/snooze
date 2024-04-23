package apiserver

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/common/health"
	"github.com/japannext/snooze/common/logging"
	"github.com/japannext/snooze/common/opensearch"
	"github.com/japannext/snooze/common/redis"
)

func Run() {
	// Init components
	logging.Init()
	health.Init()
	initConfig()
	opensearch.Init()
	redis.Init()

	// Routes
	r := gin.Default()
	r.GET("/livez", health.Heath.LiveRoute)
	r.GET("/readyz", health.Heath.ReadyRoute)
	r.Group("/static", eTagMiddleware()).Static("/", config.StaticPath)
	registerRoutes(r)

	// Listen
	hostport := fmt.Sprintf("%s:%d", config.ListenAddress, config.ListenPort)
	r.Run(hostport)
}
