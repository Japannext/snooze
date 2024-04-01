package api

import (
  "fmt"

  "github.com/gin-gonic/gin"
  log "github.com/sirupsen/logrus"

  "github.com/japannext/snooze/common/opensearch"
  "github.com/japannext/snooze/common/logging"
)

var database *opensearch.Database
var err error

// Add the API routes to the HTTP server
func mountRoutes(r *gin.Engine) {

  // Build etags map at startup

  // LogV2
  r.GET("/api/log/v2", searchLogV2)
  // Health checks
  r.GET("/livez", livez)
  r.GET("/readyz", readyz)
  // Static files
  // We wrap it in a group to add a middleware
  r.Group("/static", eTagMiddleware()).Static("/", config.StaticPath)
}

func Run() {

  if err := config.init(); err != nil {
    log.Fatal(err)
  }
  if err := logging.Init(); err != nil {
    log.Fatal(err)
  }
  log.Debugf("Loaded config: %+v", config)

  database, err = opensearch.Init()
  if err != nil {
    log.Fatal(err)
  }

  // Fail immediately if the database is unreachable
  if err := database.CheckHealth(); err != nil {
    log.Fatal(err)
  }

  r := gin.Default()
  mountRoutes(r)

  hostport := fmt.Sprintf("%s:%d", config.ListenAddress, config.ListenPort)
  r.Run(hostport)
}
