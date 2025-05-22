package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/opensearch"
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
	if err := c.ShouldBindQuery(&params); err != nil {
		c.String(http.StatusBadRequest, "invalid parameter: %s", err)
		return
	}
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
	if err := c.ShouldBindJSON(&item); err != nil {
		c.String(http.StatusBadRequest, "invalid data: %s", err)
		return
	}

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
