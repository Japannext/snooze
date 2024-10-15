package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"github.com/japannext/snooze/pkg/models"
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

func StoreLogs(ctx context.Context, index string, items []*models.Log) error {
	var buf bytes.Buffer
	indexLine := fmt.Sprintf(`{"index": {"_index": "%s"}}`, index)
	for _, item := range items {
		b, err := json.Marshal(item)
		if err != nil {
			log.Warnf("can't marshal into json `%+v`: %s", item, err)
			continue
		}
		buf.WriteString(indexLine)
		buf.WriteString("\n")
		buf.Write(b)
		buf.WriteString("\n")
	}
	req := opensearchapi.BulkReq{
		Index: index,
		Body: bytes.NewReader(buf.Bytes()),
	}
	resp, err := client.Bulk(ctx, req)
	if err != nil {
		return err
	}
	if resp.Errors {
		log.Warnf("Bulk query: %s", buf.String())
		return bulkRespToError(resp)
	}

	return nil
}

func Bulk(ctx context.Context, req opensearchapi.BulkReq) (*opensearchapi.BulkResp, error) {
	return client.Bulk(ctx, req)
}

func bulkRespToError(resp *opensearchapi.BulkResp) error {
	var buf strings.Builder
	for i, v := range resp.Items {
		for k, r := range v {
			if r.Error != nil {
				msg := fmt.Sprintf("[#%d:%s] type='%s' reason='%s'\n", i, k, r.Error.Type, r.Error.Reason)
				buf.WriteString(msg)
			}
		}
	}
	return fmt.Errorf("error in bulk log:\n%s", buf.String())
}
