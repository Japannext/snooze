package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/japannext/snooze/pkg/common/tracing"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

type UpdateReq struct {
	Index  string
	ID     string
	Params opensearchapi.UpdateParams
	Doc    map[string]interface{}
	Upsert map[string]interface{}
}

func Update(ctx context.Context, req UpdateReq) error {
	ctx, span := tracer.Start(ctx, "Update")
	defer span.End()

	body, err := json.Marshal(req.Doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document (%+v): %s", req.Doc, err)
	}
	tracing.SetString(span, "query", string(body))
	resp, err := client.Update(ctx, opensearchapi.UpdateReq{
		Index:      req.Index,
		DocumentID: req.ID,
		Params:     req.Params,
		Body:       bytes.NewReader(body),
	})
	if err != nil {
		return err
	}
	httpResp := resp.Inspect().Response
	if httpResp.IsError() {
		return fmt.Errorf("Error in update by query: %s", httpResp)
	}
	return nil
}
