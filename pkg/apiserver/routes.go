package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/pkg/common/opensearch"
)

type Controller struct {
}

func registerRoutes(r *gin.Engine) {
	r.GET("/api/alert-events/v2", searchAlertEvents)
}

func searchAlertEvents(c *gin.Context) {
	// Extract parameters
	query := c.Param("query")

	pagination := parsePagination(c)
	ll, err := opensearch.LogStore.Search(c, query, pagination)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching logs from database: %w", err)
	}

	c.JSON(http.StatusOK, ll)
}
