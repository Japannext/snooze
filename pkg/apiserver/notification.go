package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/models"
)

func getNotifications(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "getNotifications")
	defer span.End()

	var req *opensearch.SearchRequest[*models.Notification]
	req.Index = models.NOTIFICATION_INDEX

	// Pagination
	pagination := models.NewPagination()
	c.BindQuery(&pagination)
	if pagination.OrderBy == "" {
		pagination.OrderBy = "timestamp.display"
	}
	req.WithPagination(pagination)

    // Timerange
    timerange := &models.TimeRange{}
    c.BindQuery(&timerange)
	req.WithTimeRange("timestamp.display", timerange)

	// Search
    search := &models.Search{}
	c.BindQuery(&search)
	req.WithSearch(search.Text)

	items, err := req.Do(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting notification: %s", err)
		return
	}

	c.JSON(http.StatusOK, items)
}


func init() {
	routes = append(routes, func(r *gin.Engine) {
		r.GET("/api/notifications", getNotifications)
	})
}
