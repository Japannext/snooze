package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/pkg/common/opensearch"
)

func registerAlertRoutes(r *gin.Engine) {
	r.GET("/api/alert/:uid", getAlert)
	r.GET("/api/alerts")
}

func getAlert(c *gin.Context) {
	uid := c.Param("uid")

	alert, err := opensearch.LogStore.Get(uid)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting alert uid=%s: %w", uid, err)
		return
	}
	if alert == nil {
		c.String(http.StatusNotFound, "Could not find alert uid=%s", uid)
		return
	}

	c.JSON(http.StatusOK, alert)
}
