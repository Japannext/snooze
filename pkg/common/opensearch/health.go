package opensearch

import (
	"context"
	"fmt"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

// Check if the session is up.
func CheckHealth(ctx context.Context) error {
	resp, err := client.Cluster.Health(ctx, &opensearchapi.ClusterHealthReq{})
	if err != nil {
		return err
	}

	if resp.Status == "green" {
		return nil
	}
	return fmt.Errorf("Status is '%s' (not 'green')", resp.Status)
}
