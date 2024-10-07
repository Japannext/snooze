package apiserver

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	dsl "github.com/mottaquikarim/esquerydsl"

	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/models"
)

type Search struct {
	Text string `form:"search"`
}

func searchLogs(c *gin.Context) {
	var (
		start = time.Now()
		pagination = models.NewPagination()
		timerange *models.TimeRange
		search Search
	)
	c.BindQuery(&pagination)
	c.BindQuery(&timerange)
	c.BindQuery(&search)

	doc := &dsl.QueryDoc{}
	params := &opensearchapi.SearchParams{}

	if pagination.OrderBy == "" {
		pagination.OrderBy = "timestampMillis"
	}
	opensearch.AddTimeRange(doc, timerange)
	opensearch.AddPagination(doc, params, pagination)
	opensearch.AddSearch(doc, search.Text)

	items, err := opensearch.Search[*models.Log](c, models.LOG_INDEX, params, doc)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting logs for search='%s': %s", search.Text, err)
		return
	}

	c.JSON(http.StatusOK, items)
	logSearchDuration.Observe(time.Since(start).Seconds())
}

func init() {
	routes = append(routes, func(r *gin.Engine) {
		r.GET("/api/logs", searchLogs)
	})
}
