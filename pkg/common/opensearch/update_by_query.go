package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/japannext/snooze/pkg/common/opensearch/dsl"
	"github.com/japannext/snooze/pkg/common/tracing"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

type UpdateByQueryReq struct {
	Index  string
	Params opensearchapi.UpdateByQueryParams
	Doc    dsl.QueryReq
}

func UpdateByQuery(ctx context.Context, req UpdateByQueryReq) error {
	ctx, span := tracer.Start(ctx, "UpdateByQuery")
	defer span.End()

	body, err := json.Marshal(req.Doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document (%+v): %s", req.Doc, err)
	}
	tracing.SetString(span, "query", string(body))
	resp, err := client.UpdateByQuery(ctx, opensearchapi.UpdateByQueryReq{
		Indices: []string{req.Index},
		Params:  req.Params,
		Body:    bytes.NewReader(body),
	})
	if err != nil {
		return err
	}
	tracing.SetInt(span, "response.updated", resp.Updated)
	tracing.SetInt(span, "response.total", resp.Total)
	tracing.SetInt(span, "response.failures", len(resp.Failures))
	httpResp := resp.Inspect().Response
	if httpResp.IsError() {
		return fmt.Errorf("Error in update by query: %s", httpResp)
	}
	return nil
}
