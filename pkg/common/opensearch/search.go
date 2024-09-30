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

func Search[T api.HasID](ctx context.Context, index string, params *opensearchapi.SearchParams, doc *dsl.QueryDoc) (*api.ListOf[T], error) {
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
		(*list.Items[i]).SetID(hit.ID)
	}
	list.Total = resp.Hits.Total.Value
	if resp.Hits.Total.Relation != "eq" {
		list.More = true
	}
    return &list, nil
}
