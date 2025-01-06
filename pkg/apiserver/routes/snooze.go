package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/opensearch/dsl"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/utils"
	"github.com/japannext/snooze/pkg/models"
)

func registerSnoozeRoutes(r *gin.Engine) {
	r.GET("/snoozes", getSnoozes)
	r.POST("/snooze", postSnooze)
	r.POST("/snooze/cancel", postSnoozeCancel)
}

type getSnoozesParams struct {
	*models.Pagination
	*models.TimeRange
	*models.Search
	*models.Filter
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

	if params.Filter != nil {
		now := uint64(time.Now().UnixMilli())
		switch params.Filter.Text {
		case "active", "":
			req.Doc.Or(dsl.QueryItem{Range: map[string]dsl.Range{"startAt": {Lte: &now}}})
			req.Doc.Or(dsl.QueryItem{Range: map[string]dsl.Range{"expireAt": {Gte: &now}}})
			req.Doc.WithNotExists("cancelled.reason")
		case "expired":
			req.Doc.WithRange("expireAt", dsl.Range{Lt: &now})
			req.Doc.WithNotExists("cancelled.reason")
		case "upcoming":
			req.Doc.WithRange("startAt", dsl.Range{Gt: &now})
			req.Doc.WithNotExists("cancelled.reason")
		case "cancelled":
			req.Doc.WithExists("cancelled.reason")
		case "all":
			// no filter
		default:
			c.String(http.StatusBadRequest, "unknown filter name `%s`", params.Filter.Text)
			return
		}
	}

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

	var item *models.Snooze
	err := c.BindJSON(&item)
	if err != nil {
		c.String(http.StatusBadRequest, "error parsing snooze object: %s", err)
		return
	}

	item.Username = c.GetString("username")

	data, err := json.Marshal(item)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to marshal Snooze entry: %s", err)
		return
	}

	// Compute hash if not already computed
	for _, group := range item.Groups {
		if group.Hash == "" {
			group.Hash = utils.ComputeHash(group.Labels)
		}
	}

	pipe := redis.Client.Pipeline()
	for _, group := range item.Groups {
		key := fmt.Sprintf("snooze/%s/%s", group.Name, group.Hash)
		pipe.Set(ctx, key, data, 0)
		pipe.ExpireAt(ctx, key, item.EndsAt.Time)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		c.String(http.StatusInternalServerError, "failed to execute redis pipeline: %s", err)
		return
	}

	err = opensearch.Index(ctx, &opensearch.IndexReq{
		Index: models.SNOOZE_INDEX,
		Item:  item,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to create snooze entry: %s", err)
		return
	}
}

type postSnoozeCancelParams struct {
	// The snooze entries to cancel. Must at least contain the `_id` and `groups`
	IDs []string `json:"ids"`
	// The reason why it was cancelled
	Reason string `json:"reason"`
}

func postSnoozeCancel(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "postSnoozeCancel")
	defer span.End()

	params := postSnoozeCancelParams{}
	c.BindJSON(&params)

	if len(params.IDs) == 0 {
		c.String(http.StatusBadRequest, "no snooze ID specified in `ids`")
		return
	}

	searchReq := &opensearch.SearchReq{
		Index: models.SNOOZE_INDEX,
	}
	searchReq.Doc.WithTerms("_id", params.IDs)

	list, err := opensearch.Search[*models.Snooze](ctx, searchReq)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to find snooze entries: %s", err)
		return
	}

	var keys []string
	for _, item := range list.Items {
		for _, group := range item.Groups {
			keys = append(keys, fmt.Sprintf("snooze/%s/%s", group.Name, group.Hash))
		}
	}
	err = redis.Client.Del(ctx, keys...).Err()
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to delete keys from redis: %s", err)
		return
	}

	updateReq := opensearch.UpdateByQueryReq{
		Index: models.SNOOZE_INDEX,
	}
	updateReq.Doc.WithTerms("_id", params.IDs)
	script := strings.Join([]string{
		"ctx._source.cancelled = params.emptyobject",
		"ctx._source.cancelled.reason = params.reason",
		"ctx._source.cancelled.by = params.by",
		"ctx._source.cancelled.at = params.at",
	}, ";")
	updateReq.Doc.WithPainlessScript(script, map[string]interface{}{
		"emptyobject": map[string]string{},
		"reason":      params.Reason,
		"at":          time.Now().UnixMilli(),
		"by":          "",
	})

	err = opensearch.UpdateByQuery(ctx, updateReq)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to update snooze in opensearch: %s", err)
		return
	}
}
