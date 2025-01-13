package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/models"
)

type Filter struct {
	Text string `form:"filter"`
}

type getLogsParams struct {
	*models.Pagination
	*models.TimeRange
	*models.Search
	*Filter
}

func getLogs(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "getLogs")
	defer span.End()

	start := time.Now()

	req := &opensearch.SearchReq{Index: models.LogIndex}

	params := getLogsParams{Pagination: models.NewPagination()}
	c.BindQuery(&params)
	if params.Pagination.OrderBy == "" {
		params.Pagination.OrderBy = "displayTime"
	}
	req.WithPagination(params.Pagination)
	req.WithTimeRange("displayTime", params.TimeRange)
	req.WithSearch(params.Search)

	// Filters
	if params.Filter != nil {
		switch params.Filter.Text {
		case "active", "":
			req.Doc.WithTerm("status.kind.keyword", models.LogActive)
		case "snoozed":
			req.Doc.WithTerm("status.kind.keyword", models.LogSnoozed)
		case "acked":
			req.Doc.WithTerm("status.kind.keyword", models.LogAcked)
		case "failed":
			// req.Doc.WithTerm("", "")
		case "all":
			// no filter
		default:
			c.String(http.StatusBadRequest, "unknown filter name `%s`", params.Filter.Text)
			return
		}
	}

	items, err := opensearch.Search[*models.Log](ctx, req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting log: %s", err)
		return
	}

	c.JSON(http.StatusOK, items)
	logSearchDuration.Observe(time.Since(start).Seconds())
}
