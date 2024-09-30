package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	dsl "github.com/mottaquikarim/esquerydsl"

	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/redis"
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func listAlerts(c *gin.Context) {

	var pagination = api.NewPagination()
	c.BindQuery(&pagination)

	doc := &dsl.QueryDoc{}
	params := &opensearchapi.SearchParams{}

	if pagination.OrderBy == "" {
		pagination.OrderBy = "startsAt"
	}
	// opensearch.AddTimeRange(doc, timerange)
	opensearch.AddPagination(doc, params, pagination)

	res, err := opensearch.Search[*api.Alert](c, api.ALERT_INDEX, params, doc)
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
