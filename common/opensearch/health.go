package opensearch

import (
  "bytes"
  "context"
  "encoding/json"
  "fmt"

  api "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

type ClusterHealth struct {
  ClusterName string `json:"cluster_name"`
  Status string `json:"status"`
}

// Check if the session is up
func (client *OpensearchClient) CheckHealth() error {
  ctx := context.Background()
  resp, err := api.ClusterHealthRequest{}.Do(ctx, client.Client)
  if err != nil {
    return err
  }
  ch := &ClusterHealth{}

  var buf bytes.Buffer
  buf.ReadFrom(resp.Body)
  err = json.Unmarshal(buf.Bytes(), &ch)
  if err != nil {
    return fmt.Errorf("Can't unmarshal cluster health: %w", err)
  }
  if ch.Status == "green" {
    return nil
  }
  return fmt.Errorf("Status is '%s' (not 'green')", ch.Status)
}
