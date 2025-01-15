package opensearch

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"

    "github.com/japannext/snooze/pkg/common/opensearch/dsl"
    "github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

type DeleteByQueryReq struct {
	Index string
	Doc dsl.QueryReq
}

func DeleteByQuery(ctx context.Context, req DeleteByQueryReq) error {
	body, err := json.Marshal(req.Doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document (%+v): %s", req.Doc, err)
	}

	resp, err := client.Document.DeleteByQuery(ctx, opensearchapi.DocumentDeleteByQueryReq{
		Indices: []string{req.Index},
		Body: bytes.NewReader(body),
	})
	if err != nil {
		return fmt.Errorf("failed to delete_by_query: %w", err)
	}

	httpResp := resp.Inspect().Response
	if httpResp.IsError() {
		return fmt.Errorf("error in delete by query: %w", err)
	}

	return nil
}
