package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/opensearch"
)

type getNotificationsParams struct {
	*models.Pagination
	*models.TimeRange
	*models.Search
}

func getNotifications(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "getNotifications")
	defer span.End()

	req := &opensearch.SearchReq{Index: models.NotificationIndex}

	// Params
	params := getNotificationsParams{Pagination: models.NewPagination()}
	if err := c.ShouldBindQuery(&params); err != nil {
		c.String(http.StatusBadRequest, "invalid parameter: %s", err)
		return
	}
	if params.OrderBy == "" {
		params.OrderBy = "notificationTime"
	}
	req.WithPagination(params.Pagination)
	req.WithTimeRange("notificationTime", params.TimeRange)
	req.WithSearch(params.Search)

	items, err := opensearch.Search[*models.Notification](ctx, req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting notification: %s", err)
		return
	}

	c.JSON(http.StatusOK, items)
}
