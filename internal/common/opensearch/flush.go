package opensearch

import (
	"context"
)

func Flush(ctx context.Context, index string) error {
	req := DeleteByQueryReq{Index: index}
	req.Doc.MatchAll()

	return DeleteByQuery(ctx, req)
}
