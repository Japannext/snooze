package apiserver

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/models"
)

func init() {
	routes = append(routes, func(r *gin.Engine) {
		r.GET("/api/logs", getLogs)
	})
}

type getLogsParams struct {
	*models.Pagination
	*models.TimeRange
	*models.Search
	filter string `form:"filter"`
}

func getLogs(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "getLogs")
	defer span.End()

	start := time.Now()

	req := &opensearch.SearchReq{Index: models.LOG_INDEX}

	params := getLogsParams{Pagination: models.NewPagination()}
	c.BindQuery(&params)
	if params.Pagination.OrderBy == "" {
		params.Pagination.OrderBy = "displayTime"
	}
	req.WithPagination(params.Pagination)
	req.WithTimeRange("displayTime", params.TimeRange)
	req.WithSearch(params.Search)

	// Log filters
	switch params.filter {
	case "active", "":
		// No filter
	case "snoozed":
		req.Doc.WithExists("status.snoozed")
	case "acked":
		req.Doc.WithExists("status.acked")
	case "failed":
		// req.Doc.WithTerm("", "")
	default:
		c.String(http.StatusBadRequest, "unknown filter name `%s`", params.filter)
		return
	}

	items, err := opensearch.Search[*models.Log](ctx, req)
    if err != nil {
        c.String(http.StatusInternalServerError, "Error getting log: %s", err)
        return
    }

    c.JSON(http.StatusOK, items)
	logSearchDuration.Observe(time.Since(start).Seconds())
}
