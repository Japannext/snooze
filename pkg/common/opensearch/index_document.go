package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/japannext/snooze/pkg/models"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

type IndexReq struct {
	Index  string
	ID     string
	Params opensearchapi.IndexParams
	Item   models.HasID
}

func Index(ctx context.Context, req *IndexReq) error {
	body, err := json.Marshal(req.Item)
	if err != nil {
		return fmt.Errorf("failed to marshal document (%+v): %s", req.Item, err)
	}

	resp, err := client.Index(ctx, opensearchapi.IndexReq{
		Index:      req.Index,
		DocumentID: req.ID,
		Params:     req.Params,
		Body:       bytes.NewReader(body),
	})
	if err != nil {
		return fmt.Errorf("error indexing document to '%s': %s", req.Index, err)
	}
	r := resp.Inspect().Response
	if r.IsError() {
		return fmt.Errorf("error indexing document to '%s': %s", req.Index, r.String())
	}

	return nil
}
