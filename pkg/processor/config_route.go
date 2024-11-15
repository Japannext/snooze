package processor

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	routes = append(routes, func(r *gin.Engine) {
		r.GET("/process/config", getProcessConfig)
	})
}

func getProcessConfig(c *gin.Context) {
	c.JSON(http.StatusOK, pipeline)
}
