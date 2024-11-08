package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"github.com/japannext/snooze/pkg/common/opensearch/dsl"
)

type UpdateByQueryReq struct {
	Index string
	Params opensearchapi.UpdateByQueryParams
	Doc dsl.QueryReq
}

func UpdateByQuery(ctx context.Context, req UpdateByQueryReq) (*opensearchapi.UpdateByQueryResp, error) {
	body, err := json.Marshal(req.Doc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal document (%+v): %s", req.Doc, err)
	}
	resp, err := client.UpdateByQuery(ctx, opensearchapi.UpdateByQueryReq{
		Indices: []string{req.Index},
		Params: req.Params,
		Body: bytes.NewReader(body),
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
