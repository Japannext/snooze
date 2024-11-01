package apiserver

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/models"
)

func getLogs(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "getLogs")
	defer span.End()

	start := time.Now()

    var req *opensearch.SearchRequest[*models.Log]
    req.Index = models.LOG_INDEX

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
        c.String(http.StatusInternalServerError, "Error getting log: %s", err)
        return
    }

    c.JSON(http.StatusOK, items)
	logSearchDuration.Observe(time.Since(start).Seconds())
}

func init() {
	routes = append(routes, func(r *gin.Engine) {
		r.GET("/api/logs", getLogs)
	})
}
