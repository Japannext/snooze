package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/pkg/common/opensearch"
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

type LogsResponse struct {
	Logs []api.Log `json:"logs"`
	Total int `json:"total"`
}

func registerLogRoutes(r *gin.Engine) {
	r.GET("/api/log/:uid", getLog)
	r.GET("/api/logs", searchLogs)
}

func getLog(c *gin.Context) {
	uid := c.Param("uid")

	item, err := opensearch.LogStore.GetLog(uid)
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

type Search struct {
	Text string `form:"search"`
}

func searchLogs(c *gin.Context) {

	var (
		pagination = api.NewPagination()
		timerange api.TimeRange
		search Search
	)
	c.BindQuery(&pagination)
	c.BindQuery(&timerange)
	c.BindQuery(&search)

	res, err := opensearch.LogStore.SearchLogs(c, search.Text, timerange, pagination)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting logs for search='%s': %s", search.Text, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
