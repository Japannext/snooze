package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/japannext/snooze/pkg/samples"
)

func registerSampleRoute(r *gin.Engine) {
	r.POST("/api/sample", generateSamples)
}

func generateSamples(c *gin.Context) {

	err := samples.Run()
	if err != nil {
		c.String(http.StatusInternalServerError, "Issue generating samples: %s", err)
		return
	}

	c.String(http.StatusOK, "All samples generated successfully")
}

