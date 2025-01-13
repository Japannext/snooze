package routes

import (
	"net/http"
	// "time"

	"github.com/gin-gonic/gin"
	"github.com/japannext/snooze/pkg/common/opensearch"
	// "github.com/japannext/snooze/pkg/common/opensearch/dsl"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/models"
)

const (
	alertFilterActive  = "active"
	alertFilterHistory = "history"
)

type getAlertsParams struct {
	*models.Pagination
	*models.TimeRange
	*models.Search
	*models.Filter
}

func getAlerts(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "getAlerts")
	defer span.End()

	params := getAlertsParams{Pagination: models.NewPagination()}
	if err := c.BindQuery(&params); err != nil {
		c.String(http.StatusBadRequest, "error in query parameter: %s", err)

		return
	}

	if params.Pagination.OrderBy == "" {
		params.Pagination.OrderBy = "startsAt"
	}

	req := &opensearch.SearchReq{}

	if params.Filter != nil {
		switch params.Filter.Text {
		case alertFilterActive:
			req.Index = models.ActiveAlertIndex
		case alertFilterHistory:
			req.Index = models.AlertHistoryIndex
		default:
			c.String(http.StatusBadRequest, "unknown filter name `%s`", params.Filter.Text)
			return
		}
	}

	req.WithPagination(params.Pagination)
	req.WithTimeRange("startsAt", params.TimeRange)
	req.WithSearch(params.Search)

	items, err := opensearch.Search[*models.AlertRecord](ctx, req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting alerts for : %s", err)
		return
	}

	c.JSON(http.StatusOK, items)
}

func postAlert(c *gin.Context) {
	// TODO
}

type AlertKeys struct {
	Keys []string
}

func getLiveStatus(c *gin.Context) {
	var keys []string
	c.BindJSON(&keys)

	items := make([]string, len(keys))

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
