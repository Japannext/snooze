package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	// dsl "github.com/mottaquikarim/esquerydsl"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/opensearch/dsl"
)

type SearchRequest[T models.HasID] struct {
	Index string
	Params opensearchapi.SearchParams
	Doc dsl.QueryRequest
}

func (req *SearchRequest[T]) Do(ctx context.Context) (*models.ListOf[T], error) {

    body, err := json.Marshal(req.Doc)
    if err != nil {
        return nil, fmt.Errorf("invalid request body (%+v): %w", req.Doc, err)
    }
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
	list.Items = make([]*T, len(resp.Hits.Hits))
	for i, hit := range resp.Hits.Hits {
		if err := json.Unmarshal(hit.Source, &list.Items[i]); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message %s: %w", hit.Source, err)
		}
		(*list.Items[i]).SetID(hit.ID)
	}
	list.Total = resp.Hits.Total.Value
	if resp.Hits.Total.Relation != "eq" {
		list.More = true
	}
    return &list, nil
}

func (req *SearchRequest[T]) WithTerm(field, value string) *SearchRequest[T] {
	item := dsl.QueryItem{Term: kv(field, value)}
	req.Doc.Query.Bool.And = append(req.Doc.Query.Bool.And, item)
	return req
}

func (req *SearchRequest[T]) WithSort(field string, ascending bool) *SearchRequest[T] {
	var order string
	if ascending {
		order = "asc"
	} else {
		order = "desc"
	}
	req.Doc.Sort = append(req.Doc.Sort, kv(field, dsl.Sort{Order: order}))
	return req
}

// A wrapper to wrap all the single-value maps used by opensearch
func kv[T any](field string, value T) map[string]T {
	return map[string]T{field: value}
}

func (req *SearchRequest[T]) WithTimeRange(field string, timerange *models.TimeRange) *SearchRequest[T] {
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
		item := dsl.QueryItem{Range: kv(field, r)}
		req.Doc.Query.Bool.And = append(req.Doc.Query.Bool.And, item)
	}
	return req
}

func (req *SearchRequest[T]) WithPagination(pagination *models.Pagination) *SearchRequest[T] {
	if pagination == nil {
		return req
	}
	req.Params.From = pointer((pagination.Page - 1) * pagination.Size)
	req.Params.Size = pointer(pagination.Size)
	sort := map[string]string{}
	if pagination.Ascending {
		sort[pagination.OrderBy] = "asc"
	} else {
		sort[pagination.OrderBy] = "desc"
	}
	req.WithSort(pagination.OrderBy, pagination.Ascending)
	return req
}

func (req *SearchRequest[T]) WithSearch(text string) *SearchRequest[T] {
	if text == "" {
		return req
	}
	item := dsl.QueryItem{QueryString: &dsl.QueryString{Query: text}}
	req.Doc.Query.Bool.And = append(req.Doc.Query.Bool.And, item)
	return req
}

/*
func Search[T models.HasID](ctx context.Context, index string, params *opensearchapi.SearchParams, req.Doc *dsl.QueryDoc) (*models.ListOf[T], error) {
    body, err := json.Marshal(req.Doc)
    if err != nil {
        return nil, fmt.Errorf("invalid request body (%+v): %w", req.Doc, err)
    }
    req := &opensearchapi.SearchReq{
        Indices: []string{index},
        Params: *params,
        Body:  bytes.NewReader(body),
    }
    resp, err := client.Search(ctx, req)
    if err != nil {
        return nil, err
    }
    if resp.Errors {
        return nil, fmt.Errorf("opensearch returned an error: %s", "")
    }
	list := models.ListOf[T]{}
	list.Items = make([]*T, len(resp.Hits.Hits))
	for i, hit := range resp.Hits.Hits {
		if err := json.Unmarshal(hit.Source, &list.Items[i]); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message %s: %w", hit.Source, err)
		}
		(*list.Items[i]).SetID(hit.ID)
	}
	list.Total = resp.Hits.Total.Value
	if resp.Hits.Total.Relation != "eq" {
		list.More = true
	}
    return &list, nil
}
*/
