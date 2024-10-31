package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/opensearch/dsl"
	"github.com/japannext/snooze/pkg/models"
)

type GroupBy struct {
	Group string `form:"group"`
}

func getLogGroups(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "getLogGroups")
	defer span.End()

	var (
		groupBy = GroupBy{Group: "by-host"}
		timerange *models.TimeRange
		pagination = models.NewPagination()
	)

	c.BindQuery(&pagination)
	c.BindQuery(&timerange)
	c.BindQuery(&groupBy)

	filter := dsl.NewTermQuery("groups.name.keyword", groupBy.Group)
	doc := dsl.NewAggregationQuery("groups.hash.keyword", filter)
	params := &opensearchapi.SearchParams{}

	_, err := opensearch.Aggregate(ctx, models.LOG_INDEX, params, doc)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting groups: %s", err)
	}

	// TODO
}

func init() {
	routes = append(routes, func(r *gin.Engine) {
		r.GET("/api/groups", getLogGroups)
	})
}
