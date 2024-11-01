package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/models"
)

func getAlerts(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "getAlerts")
	defer span.End()

	var req *opensearch.SearchRequest[*models.Alert]
	req.Index = models.ALERT_INDEX

	// Pagination
	pagination := models.NewPagination()
	c.BindQuery(&pagination)
	if pagination.OrderBy == "" {
		pagination.OrderBy = "startsAt"
	}
	req.WithPagination(pagination)

    // Timerange
    timerange := &models.TimeRange{}
    c.BindQuery(&timerange)
	req.WithTimeRange("startAt", timerange)

	// Search
    search := &models.Search{}
	c.BindQuery(&search)
	req.WithSearch(search.Text)

	items, err := req.Do(ctx)
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

func init() {
	routes = append(routes, func(r *gin.Engine) {
		r.GET("/api/alerts", getAlerts)
	})
}
