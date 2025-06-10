package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/opensearch"
)

type GetAutocompleteLogParams struct {
	Kind  string `json:"kind"`
	Key   string `json:"key"`
	Value string `json:"value"`
	// Field  pipeline.Field `json:"field"`
}

// Return the autocomplete of a log field
func getAutocompleteLog(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "getAutocompleteLog")
	defer span.End()

	params := GetAutocompleteLogParams{}
	err := c.ShouldBind(&params)
	if err != nil {
		c.String(http.StatusBadRequest, "wrong field provided: %s", err)
		return
	}

	req := &opensearch.SearchReq{Index: models.LogIndex}
	opensearch.Search[*models.Log](ctx, req)

	switch params.Kind {
	case "label":
		req.Doc.WithTerm(fmt.Sprintf("labels.%s", params.Key), params.Value)
	case "identity":
		req.Doc.WithTerm(fmt.Sprintf("identity.%s", params.Key), params.Value)
	default:
		c.String(http.StatusBadRequest, "param kind=%s not supported", params.Kind)
		return
	}

	req.WithSize(10)

	list, err := opensearch.Search[*models.Log](ctx, req)
	if err != nil {

		return
	}

	results := []string{}
	for _, item := range list.Items {
		switch params.Kind {
		case "label":
			value, found := item.Labels[params.Key]
			if !found {
				continue
			}
			results = append(results, value)
		}
	}

	c.JSON(http.StatusOK, results)
}
