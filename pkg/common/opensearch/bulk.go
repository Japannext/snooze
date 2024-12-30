package opensearch

import (
	"context"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

func Bulk(ctx context.Context, req opensearchapi.BulkReq) (*opensearchapi.BulkResp, error) {
	return client.Bulk(ctx, req)
}
