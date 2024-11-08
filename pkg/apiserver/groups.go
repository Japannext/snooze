package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/tracing"
	"github.com/japannext/snooze/pkg/models"
)

func init() {
	routes = append(routes, func(r *gin.Engine) {
		r.GET("/api/groups", getGroups)
	})
}

type getGroupParams struct {
	*models.Search
}

func getGroups(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "getGroups")
	defer span.End()

	req := &opensearch.SearchReq{Index: models.GROUP_INDEX}

	// Search related to group
	params := getGroupParams{}
	c.BindQuery(&params)
	req.WithSearch(params.Search)

    items, err := opensearch.Search[*models.Group](ctx, req)
    if err != nil {
        c.String(http.StatusInternalServerError, "Error getting log: %s", err)
		tracing.Error(span, err)
        return
    }

    c.JSON(http.StatusOK, items)
}
