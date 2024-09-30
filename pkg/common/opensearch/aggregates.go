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

type AggregationResult struct {
	DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
	SumOtherDocCount int `json:"sum_other_doc_count"`
	Buckets []AggregationBucket `json:"buckets"`
}

type AggregationBucket struct {
	Key string `json:"key"`
	DocCount int `json:"doc_count"`
}

func Aggregate[T api.HasID](ctx context.Context, index string, params *opensearchapi.SearchParams, doc *dsl.QueryDoc) (map[string]AggregationResult, error) {
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
	if len(resp.Aggregations) == 0 {
		// TODO
	}
	var agg map[string]AggregationResult
	if err := json.Unmarshal(resp.Aggregations, &agg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message %s: %w", resp.Aggregations, err)
	}
	return agg, nil
}

