package apiserver

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/pkg/common/opensearch"
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

type Search struct {
	Text string `form:"search"`
}

func searchLogs(c *gin.Context) {

	var (
		start = time.Now()
		pagination = api.NewPagination()
		timerange *api.TimeRange
		search Search
	)
	c.BindQuery(&pagination)
	c.BindQuery(&timerange)
	c.BindQuery(&search)

	res, err := opensearch.SearchLogs(c, search.Text, timerange, pagination)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting logs for search='%s': %s", search.Text, err)
		return
	}

	c.JSON(http.StatusOK, res)
	logSearchDuration.Observe(time.Since(start).Seconds())
}

func init() {
	routes = append(routes, func(r *gin.Engine) {
		r.GET("/api/logs", searchLogs)
	})
}
