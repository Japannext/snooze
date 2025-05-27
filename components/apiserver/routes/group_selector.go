package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/opensearch"
	"github.com/japannext/snooze/internal/tracing"
)

type getGroupSelectorsParams struct {
	*models.Search
}

func getGroupSelectors(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "getGroupSelectors")
	defer span.End()

	req := &opensearch.SearchReq{}
	req.WithSize(1000)
	req.Doc.WithSort("score", false)
	req.Doc.WithSort("name", true)

	// Search related to group
	params := getGroupsParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		c.String(http.StatusBadRequest, "invalid parameter: %s", err)

		return
	}

	req.WithSearch(params.Search)

	items, err := opensearch.Search[*models.Group](ctx, req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting log: %s", err)
		tracing.Error(span, err)
		data, _ := json.Marshal(req.Doc)
		tracing.SetString(span, "opensearch.query", string(data))
		return
	}

	c.JSON(http.StatusOK, items)
}
