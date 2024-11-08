package apiserver

import (
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/models"
)

func init() {
	routes = append(routes, func(r *gin.Engine) {
		r.GET("/api/snoozes", getSnoozes)
		r.POST("/api/snooze", postSnooze)
	})
}

type getSnoozesParams struct {
	*models.Pagination
	*models.TimeRange
	*models.Search
}

func getSnoozes(c *gin.Context) {
    ctx, span := tracer.Start(c.Request.Context(), "getSnoozes")
    defer span.End()

	req := &opensearch.SearchReq{Index: models.SNOOZE_INDEX}

	// Params
	params := getSnoozesParams{Pagination: models.NewPagination()}
	c.BindQuery(&params)
	if params.Pagination.OrderBy == "" {
		params.Pagination.OrderBy = "startAt"
	}
	req.WithPagination(params.Pagination)
	req.WithTimeRange("startAt", params.TimeRange)
	req.WithSearch(params.Search)

    items, err := opensearch.Search[*models.Snooze](ctx, req)
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

	data, err := json.Marshal(item)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to marshal Snooze entry: %s", err)
		return
	}

	pipe := redis.Client.Pipeline()
	for _, group := range item.Groups {
		key := fmt.Sprintf("snooze/%s/%s", group.Name, group.Hash)
		pipe.Set(ctx, key, data, 0)
		pipe.ExpireAt(ctx, key, item.ExpireAt.Time)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		c.String(http.StatusInternalServerError, "failed to execute redis pipeline: %s", err)
		return
	}

	if err := storeQ.PublishData(ctx, opensearch.Create(models.SNOOZE_INDEX, item)); err != nil {
		c.String(http.StatusInternalServerError, "failed to publish snooze entry: %s", err)
		return
	}
}
