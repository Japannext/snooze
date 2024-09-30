package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/redis"
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func listAlerts(c *gin.Context) {

	var history = false
	var pagination = api.NewPagination()
	c.BindQuery(&pagination)

	res, err := opensearch.SearchAlerts(c, pagination, history)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting alerts for : %s", err)
		return
	}

	c.JSON(http.StatusOK, res)
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
		r.GET("/api/alerts", listAlerts)
	})
}
