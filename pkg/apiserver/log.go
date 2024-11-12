package apiserver

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/models"
)

func init() {
	routes = append(routes, func(r *gin.Engine) {
		r.GET("/api/logs", getLogs)
	})
}

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

	req := &opensearch.SearchReq{Index: models.LOG_INDEX}

	params := getLogsParams{Pagination: models.NewPagination()}
	c.BindQuery(&params)
	if params.Pagination.OrderBy == "" {
		params.Pagination.OrderBy = "displayTime"
	}
	req.WithPagination(params.Pagination)
	req.WithTimeRange("displayTime", params.TimeRange)
	req.WithSearch(params.Search)

	if params.Filter != nil {
		log.Debugf("params.filter: %s", params.Filter.Text)
	}
	if params.Search != nil {
		log.Debugf("params.search: %s", params.Search.Text)
	}
	log.Debugf("params.timerange.start: %d", params.TimeRange.Start)
	log.Debugf("params.timerange.end: %d", params.TimeRange.End)
	log.Debugf("params.pagination.page: %d", params.Pagination.Page)
	log.Debugf("params.pagination.size: %d", params.Pagination.Size)

	// Filters
	if params.Filter != nil {
		switch params.Filter.Text {
		case "active", "":
			req.Doc.WithTerm("status.kind.keyword", "active")
		case "snoozed":
			req.Doc.WithTerm("status.kind.keyword", "snoozed")
		case "acked":
			req.Doc.WithTerm("status.kind.keyword", "acked")
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
