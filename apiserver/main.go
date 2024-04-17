package apiserver

import (
  "fmt"

  "github.com/gin-gonic/gin"
  log "github.com/sirupsen/logrus"

  "github.com/japannext/snooze/common/opensearch"
  "github.com/japannext/snooze/common/logging"
)

func Run() {
  // Init components
  health.Init()
  logging.Init()
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
