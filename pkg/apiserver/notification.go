package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/pkg/common/opensearch"
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

type NotificationsResponse struct {
	Notifications []api.Notification `json:"logs"`
	Total int `json:"total"`
}

/*
func getNotification(c *gin.Context) {
	uid := c.Param("uid")

	item, err := opensearch.LogStore.GetNotification(uid)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting log uid=%s: %w", uid, err)
		return
	}
	if item == nil {
		c.String(http.StatusNotFound, "Could not find log uid=%s", uid)
		return
	}

	c.JSON(http.StatusOK, item)
}
*/

func searchNotifications(c *gin.Context) {

	var (
		pagination = api.NewPagination()
		timerange api.TimeRange
		search Search
	)
	c.BindQuery(&pagination)
	c.BindQuery(&timerange)
	c.BindQuery(&search)

	res, err := opensearch.LogStore.SearchNotifications(c, search.Text, timerange, pagination)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting logs for search='%s': %s", search.Text, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func init() {
	routes = append(routes, func(r *gin.Engine) {
		r.GET("/api/notifications", searchNotifications)
	})
}
