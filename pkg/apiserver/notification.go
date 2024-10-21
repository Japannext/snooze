package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	dsl "github.com/mottaquikarim/esquerydsl"

	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/models"
)

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
		pagination = models.NewPagination()
		timerange *models.TimeRange
		search Search
	)
	c.BindQuery(&pagination)
	c.BindQuery(&timerange)
	c.BindQuery(&search)

	var params = &opensearchapi.SearchParams{}
	var doc = &dsl.QueryDoc{}
	if pagination.OrderBy == "" {
		pagination.OrderBy = "timestamp.display"
	}
	opensearch.AddTimeRange(doc, timerange)
	opensearch.AddPagination(doc, params, pagination)
	// opensearch.AddSearch(doc, text)

	res, err := opensearch.Search[*models.Notification](c, models.NOTIFICATION_INDEX, params, doc)
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
