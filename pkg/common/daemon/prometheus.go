package daemon

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (dm *DaemonManager) setupPrometheus() {
	r := dm.getGinEngine()
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
