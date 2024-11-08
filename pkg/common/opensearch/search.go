package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/opensearch/dsl"
)

type SearchReq struct {
	Index string
	Params opensearchapi.SearchParams
	Doc dsl.QueryReq
}

func (req *SearchReq) WithTimeRange(field string, timerange *models.TimeRange) *SearchReq {
	if timerange == nil {
		return req
	}
	var r = dsl.Range{}
	if timerange.Start > 0 {
		r.Gte = &timerange.Start
	}
	if timerange.End > 0 {
		r.Lte = &timerange.End
	}
	if timerange.Start > 0 || timerange.End > 0 {
		item := dsl.QueryItem{Range: map[string]dsl.Range{field: r}}
		req.Doc.And(item)
	}
	return req
}

func (req *SearchReq) WithPagination(pagination *models.Pagination) *SearchReq {
	if pagination == nil {
		pagination = models.NewPagination()
	}
	req.Params.From = pointer((pagination.Page - 1) * pagination.Size)
	req.Params.Size = pointer(pagination.Size)
	req.Doc.WithSort(pagination.OrderBy, pagination.Ascending)
	return req
}

func (req *SearchReq) WithSearch(s *models.Search) *SearchReq {
	if s == nil || s.Text == "" {
		return req
	}
	req.Doc.WithQueryString(s.Text)
	return req
}

func Search[T models.HasID](ctx context.Context, req *SearchReq) (*models.ListOf[T], error) {
    body, err := json.Marshal(req.Doc)
    if err != nil {
        return nil, fmt.Errorf("invalid request body (%+v): %w", req.Doc, err)
    }
	log.Debugf("req body: %s", body)
	log.Debugf("req params: %+v", req.Params)
    resp, err := client.Search(ctx, &opensearchapi.SearchReq{
        Indices: []string{req.Index},
        Params: req.Params,
        Body:  bytes.NewReader(body),
	})
    if err != nil {
        return nil, err
    }
    if resp.Errors {
        return nil, fmt.Errorf("opensearch returned an error: %s", "")
    }
	list := models.ListOf[T]{}
	list.Items = make([]T, len(resp.Hits.Hits))
	for i, hit := range resp.Hits.Hits {
		if err := json.Unmarshal(hit.Source, &list.Items[i]); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message %s: %w", hit.Source, err)
		}
		list.Items[i].SetID(hit.ID)
	}
	list.Total = resp.Hits.Total.Value
	if resp.Hits.Total.Relation != "eq" {
		list.More = true
	}
    return &list, nil
}

/*
func (req *SearchRequest[T]) WithTerm(field, value string) *SearchRequest[T] {
	item := dsl.QueryItem{Term: kv(field, value)}
	req.Doc.Query.Bool.And = append(req.Doc.Query.Bool.And, item)
	return req
}
*/

/*
// A wrapper to wrap all the single-value maps used by opensearch
func kv[T any](field string, value T) map[string]T {
	return map[string]T{field: value}
}
*/

/*
func (req *SearchReq) WithGroupSearch(groupBy, search string) *SearchReq {
	item := dsl.QueryItem{QueryString: &dsl.QueryString{Query: ""}}
	req.Doc.Query.Bool.And = append(req.Doc.Query.Bool.And, item)

	return req
}
*/
