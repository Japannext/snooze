package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	dsl "github.com/mottaquikarim/esquerydsl"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	api "github.com/japannext/snooze/pkg/common/api/v2"
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


// Store one object (json serializable) into an index
func store(ctx context.Context, index string, item interface{}) (string, error) {
    b, err := json.Marshal(item)
    if err != nil {
        return "", err
    }
    body := bytes.NewReader(b)
    req := opensearchapi.IndexReq{
        Index: index,
        Body:  body,
    }
    resp, err := client.Index(ctx, req)
    if err != nil {
        return "", err
    }
    if resp.Shards.Successful == 0 {
        return "", fmt.Errorf("Failed to index object to '%s': %s", index, resp.Result)
    }
    log.Debugf("inserted into `%s`: %s", index, b)
    return resp.ID, nil
}

func search[T api.HasID](ctx context.Context, index, search string, timerange api.TimeRange, pagination api.Pagination) (*api.ListOf[T], error) {
    var doc = &dsl.QueryDoc{}
    var params = &opensearchapi.SearchParams{}

    addTimeRange(doc, timerange)
    addPagination(doc, params, pagination)
    addSearch(doc, search)

    body, err := json.Marshal(doc)
    if err != nil {
        return nil, fmt.Errorf("invalid request body (%+v): %w", doc, err)
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
	list := api.ListOf[T]{}
	list.Items = make([]*T, len(resp.Hits.Hits))
	for i, hit := range resp.Hits.Hits {
		if err := json.Unmarshal(hit.Source, &list.Items[i]); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message %s: %w", hit.Source, err)
		}
	}
	list.Total = resp.Hits.Total.Value
	if resp.Hits.Total.Relation != "eq" {
		list.More = true
	}
    return &list, nil
}

type Range struct {
	Gte *uint64 `json:"gte,omitempty"`
	Lte *uint64 `json:"lte,omitempty"`
}

func addTimeRange(doc *dsl.QueryDoc, timerange api.TimeRange) {
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

func addPagination(doc *dsl.QueryDoc, params *opensearchapi.SearchParams, pagination api.Pagination) {
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

func addSearch(doc *dsl.QueryDoc, search string) {
	if search != "" {
		item := dsl.QueryItem{Type: dsl.Match, Field: "message", Value: search}
		doc.And = append(doc.And, item)
	}
}

