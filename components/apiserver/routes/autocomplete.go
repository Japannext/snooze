package routes

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/internal/lang"
	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/opensearch"
	"github.com/japannext/snooze/internal/pipeline"
)

type GetAutocompleteLogData struct {
	Search string         `json:"search"`
	Field  pipeline.Field `json:"field"`
}

// Return the autocomplete of a log field
func getAutocompleteLog(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "getAutocompleteLog")
	defer span.End()

	data := GetAutocompleteLogData{}
	err := c.BindJSON(&data)
	if err != nil {
		c.String(http.StatusBadRequest, "wrong data provided: %s", err)
		return
	}

	field, err := lang.FieldFromNestedFields(data.Field.NestedFields)
	if err != nil {
		c.String(http.StatusBadRequest, "wrong field %s: %s", data.Field.Name, err)
		return
	}

	esField := strings.Join(data.Field.NestedFields, ".")

	req := &opensearch.SearchReq{Index: models.LogIndex}
	req.Doc.WithTerm(esField, data.Search)
	req.WithSize(10)

	list, err := opensearch.Search[*models.Log](ctx, req)
	if err != nil {

		return
	}

	results := []string{}
	for _, item := range list.Items {
		txt, err := lang.ExtractField(item, field)
		if err != nil {
			log.Warningf("failed to extract field %s: %s", field, err)
		}
		results = append(results, txt)
	}

	c.JSON(http.StatusOK, results)
}
