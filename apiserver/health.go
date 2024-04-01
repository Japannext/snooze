package api

import (
  "net/http"

  "github.com/gin-gonic/gin"
)

// Return the health of the cluster
func livez(c *gin.Context) {
  c.String(http.StatusOK, "Healthy")
}

var stopping = false

// Remove the program from the load balancer if not ready.
// Use-cases:
// * Program received a SIGTERM and is stopping.
func readyz(c *gin.Context) {
  if stopping {
    c.String(http.StatusInternalServerError, "server is stopping...")
    return
  }
}
