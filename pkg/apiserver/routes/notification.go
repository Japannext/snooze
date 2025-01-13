package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/models"
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
	c.BindQuery(&params)
	if params.Pagination.OrderBy == "" {
		params.Pagination.OrderBy = "notificationTime"
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
