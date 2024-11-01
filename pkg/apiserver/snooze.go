package apiserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/models"
)

func getSnoozes(c *gin.Context) {
    ctx, span := tracer.Start(c.Request.Context(), "getSnoozes")
    defer span.End()

    var req *opensearch.SearchRequest[*models.Snooze]
    req.Index = models.SNOOZE_INDEX

    // Pagination
    pagination := models.NewPagination()
    c.BindQuery(&pagination)
    if pagination.OrderBy == "" {
        pagination.OrderBy = "startAt"
    }
    req.WithPagination(pagination)

    // Timerange
    timerange := &models.TimeRange{}
    c.BindQuery(&timerange)
    req.WithTimeRange("startAt", timerange)

    // Search
    search := &models.Search{}
    c.BindQuery(&search)
    req.WithSearch(search.Text)

    items, err := req.Do(ctx)
    if err != nil {
        c.String(http.StatusInternalServerError, "Error getting snoozes: %s", err)
        return
    }

    c.JSON(http.StatusOK, items)
}

func postSnooze(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "postSnooze")
	defer span.End()

	var item models.Snooze
	c.BindJSON(&item)

	key := fmt.Sprintf("snooze/%s/%s", item.GroupName, item.Hash)
	pipe := redis.Client.Pipeline()
	pipe.Set(ctx, key, item.Reason, 0)
	pipe.ExpireAt(ctx, key, item.ExpireAt.Time)
	if _, err := pipe.Exec(ctx); err != nil {
	}

	if err := storeQ.PublishData(ctx, opensearch.Create(models.SNOOZE_INDEX, item)); err != nil {
		// TODO
	}
}

func init() {
	routes = append(routes, func(r *gin.Engine) {
		r.GET("/api/snoozes", getSnoozes)
		r.POST("/api/snooze", postSnooze)
	})
}
