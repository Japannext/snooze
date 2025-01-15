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
	cfg, err := LoadConfig()
	if err != nil {
		c.String(http.StatusInternalServerError, "error loading config: %s", err)

		return
	}

	c.JSON(http.StatusOK, cfg)
}
