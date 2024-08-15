package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/pkg/common/opensearch"
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

type LogsResponse struct {
	Logs []api.Log `json:"logs"`
}

func registerLogRoutes(r *gin.Engine) {
	r.GET("/api/log/:uid", getLog)
	r.GET("/api/logs", getLogs)
}

func getLog(c *gin.Context) {
	uid := c.Param("uid")

	item, err := opensearch.LogStore.Get(uid)
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

func getLogs(c *gin.Context) {
	query := c.Param("query")

	pagination := parsePagination(c)
	timerange := parseTimeRange(c)

	items, err := opensearch.LogStore.Search(c, query, timerange, pagination)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting logs for query='%s': %s", query, err)
		return
	}
	resp := LogsResponse{Logs: items}

	c.JSON(http.StatusOK, resp)
}
