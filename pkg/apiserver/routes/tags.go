package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/models"
)

type getTagsParams struct {
	*models.Search
}

func getTags(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "getTags")
	defer span.End()

	req := &opensearch.SearchReq{Index: models.TagIndex}
	req.WithSize(1000)
	// Search related to tag
	params := getTagsParams{}
	c.BindQuery(&params)
	req.WithSearch(params.Search)

	// Params
	items, err := opensearch.Search[*models.Tag](ctx, req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting snoozes: %s", err)
		return
	}

	c.JSON(http.StatusOK, items)
}

func postTag(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "postTag")
	defer span.End()

	var item *models.Tag
	c.BindJSON(&item)

	if item.Name == "" {
		c.String(http.StatusBadRequest, "tag must have a non-empty name")
		return
	}

	_, err := opensearch.Index(ctx, &opensearch.IndexReq{
		Index: models.TagIndex,
		Item:  item,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to create tag: %s", err)
		return
	}
}
