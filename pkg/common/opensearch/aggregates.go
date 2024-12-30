package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"github.com/japannext/snooze/pkg/common/opensearch/dsl"
)

type AggregationResult struct {
	DocCountErrorUpperBound int                 `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int                 `json:"sum_other_doc_count"`
	Buckets                 []AggregationBucket `json:"buckets"`
}

type AggregationBucket struct {
	Key      string `json:"key"`
	DocCount int    `json:"doc_count"`
}

func Aggregate(ctx context.Context, index string, params *opensearchapi.SearchParams, doc *dsl.AggregationRequest) (*opensearchapi.SearchResp, error) {
	body, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("invalid request body (%+v): %w", doc, err)
	}
	req := &opensearchapi.SearchReq{
		Indices: []string{index},
		Params:  *params,
		Body:    bytes.NewReader(body),
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
	/*
		var agg dsl.AggregationResponse
		if err := json.Unmarshal(resp.Aggregations, &agg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message %s: %w", resp.Aggregations, err)
		}
	*/

	return resp, nil
}
