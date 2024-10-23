package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/japannext/snooze/pkg/tester"
)

func generateSamples(c *gin.Context) {

	err := tester.Run()
	if err != nil {
		c.String(http.StatusInternalServerError, "Issue generating samples: %s", err)
		return
	}

	c.String(http.StatusOK, "All samples generated successfully")
}

func init() {
	routes = append(routes, func(r *gin.Engine) {
		r.POST("/api/sample", generateSamples)
	})
}
