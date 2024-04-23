package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/common/opensearch"
)

type Controller struct {
}

func registerRoutes(r *gin.Engine) {
	r.GET("/api/alert-events/v2", searchAlertEvents)
}

func searchAlertEvents(c *gin.Context) {
	// Extract parameters
	search := c.Param("search")
	sortBy := c.Param("sort_by")
	pp, err := extractPerPage(c)
	if err != nil {
		return
	}
	page, err := extractPage(c)
	if err != nil {
		return
	}

	ll, err := opensearch.Client.SearchAlertEvent(c, search, sortBy, page, pp)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching logs from database: %w", err)
	}

	c.JSON(http.StatusOK, ll)
}
