package opensearch

import (
	dsl "github.com/mottaquikarim/esquerydsl"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"github.com/japannext/snooze/pkg/models"
)

var (
	byTimestamp = map[string]string{"timestampMillis": "desc"}
)

func sorts(ss ...map[string]string) []map[string]string {
	return ss
}

func pointer[V any](v V) *V {
	return &v
}

type Range struct {
	Gte *uint64 `json:"gte,omitempty"`
	Lte *uint64 `json:"lte,omitempty"`
}

func AddTimeRange(doc *dsl.QueryDoc, timerange *models.TimeRange) {
	if timerange == nil {
		return
	}
	var r = Range{}
	if timerange.Start > 0 {
		r.Gte = &timerange.Start
	}
	if timerange.End > 0 {
		r.Lte = &timerange.End
	}
	if timerange.Start > 0 || timerange.End > 0 {
		item := dsl.QueryItem{Type: dsl.Range, Field: "timestampMillis", Value: r}
		doc.And = append(doc.And, item)
	}
}

func AddPagination(doc *dsl.QueryDoc, params *opensearchapi.SearchParams, pagination *models.Pagination) {
	if pagination == nil {
		return
	}
	params.From = pointer((pagination.Page - 1) * pagination.Size)
	params.Size = pointer(pagination.Size)
	sort := map[string]string{}
	if pagination.Ascending {
		sort[pagination.OrderBy] = "asc"
	} else {
		sort[pagination.OrderBy] = "desc"
	}
	doc.Sort = []map[string]string{sort}
}

func AddSearch(doc *dsl.QueryDoc, search string) {
	if search == "" {
		return
	}
	item := dsl.QueryItem{Type: dsl.Match, Field: "message", Value: search}
	doc.And = append(doc.And, item)
}
