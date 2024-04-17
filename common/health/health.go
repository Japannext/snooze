package health

import (
  //"fmt"
  "net/http"

  "github.com/gin-gonic/gin"
)

var Heath *HealthStatus

type HealthStatus struct {
  ready bool
  live bool
  liveReason string
  readyReason string
}

func (h *HealthStatus) Ready() {
  h.ready = true
  h.readyReason = "OK"
}

func (h *HealthStatus) Live() {
  h.live = true
  h.liveReason = "OK"
}

func (h *HealthStatus) NotReady(reason string) {
  h.ready = false
  h.readyReason = reason
}

func (h *HealthStatus) NotLive(reason string) {
  h.live = false
  h.liveReason = reason
}

func (h *HealthStatus) ReadyRoute(c *gin.Context) {
  var s int
  if h.ready {
    s = http.StatusOK
  } else {
    s = http.StatusInternalServerError
  }
  c.String(s, h.readyReason)
}

func (h *HealthStatus) LiveRoute(c *gin.Context) {
  var s int
  if h.live {
    s = http.StatusOK
  } else {
    s = http.StatusInternalServerError
  }
  c.String(s, h.liveReason)
}

func Init() {
  Heath = &HealthStatus{
    ready: false,
    live: false, 
    readyReason: "starting...",
    liveReason: "starting...",
  }
}

// An utility for services that don't listen on HTTP
/*
func (h *HealthStatus) Serve(host string, port int) {
  r := gin.Default()
  r.GET("/readyz", h.ReadyRoute)
  r.GET("/livez", h.LiveRoute)
  r.Run(fmt.Sprintf("%s:%d", host, port))
}
*/
