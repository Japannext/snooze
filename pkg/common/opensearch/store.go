package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

// Store one object (json serializable) into an index
func Store(ctx context.Context, index string, item interface{}) (string, error) {
    b, err := json.Marshal(item)
    if err != nil {
        return "", err
    }
    body := bytes.NewReader(b)
    req := opensearchapi.IndexReq{
        Index: index,
        Body:  body,
    }
    resp, err := client.Index(ctx, req)
    if err != nil {
        return "", err
    }
    if resp.Shards.Successful == 0 {
        return "", fmt.Errorf("Failed to index object to '%s': %s", index, resp.Result)
    }
    log.Debugf("inserted into `%s`: %s", index, b)
    return resp.ID, nil
}
