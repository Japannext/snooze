package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/models"
)

func init() {
	routes = append(routes, func(r *gin.Engine) {
		r.GET("/api/alerts", getAlerts)
	})
}

type getAlertsParams struct {
	*models.Pagination
	*models.TimeRange
	*models.Search
	filter string `form:"filter"`
}

func getAlerts(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "getAlerts")
	defer span.End()

	req := &opensearch.SearchReq{Index: models.ALERT_INDEX}

    params := getLogsParams{Pagination: models.NewPagination()}
    c.BindQuery(&params)
    if params.Pagination.OrderBy == "" {
        params.Pagination.OrderBy = "startsAt"
    }
    req.WithPagination(params.Pagination)
    req.WithTimeRange("startsAt", params.TimeRange)
    req.WithSearch(params.Search)

	switch params.filter {
	}

	items, err := opensearch.Search[*models.Alert](ctx, req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting alerts for : %s", err)
		return
	}

	c.JSON(http.StatusOK, items)
}

type AlertKeys struct {
	Keys []string
}

func getLiveStatus(c *gin.Context) {
	var keys []string
	c.BindJSON(&keys)

	var items = make([]string, len(keys))

	pipe := redis.Client.Pipeline()
	for i, key := range keys {
		item, err := pipe.Get(c, key).Result()
		if err != nil {
			items[i] = ""
		}
		items[i] = item
	}

	c.JSON(http.StatusOK, items)
}

