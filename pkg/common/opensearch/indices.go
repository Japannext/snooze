package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/japannext/snooze/pkg/models"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

func ensureIndice(ctx context.Context, name string, idx models.Indice) error {
	found, _, err := hasIndice(ctx, name)
	if err != nil {
		log.Fatalf("failed to find indice %s: %s", name, err)
	}

	if !found {
		if err := createIndice(ctx, name, idx); err != nil {
			log.Fatalf("failed to create indice %s: %s", name, err)
		}
	}

	return nil
}

func createIndice(ctx context.Context, name string, idx models.Indice) error {
	log.Debugf("Creating indice '%s'...", name)

	body, err := json.Marshal(idx)
	if err != nil {
		return fmt.Errorf("error marshaling indice '%s': %w", name, err)
	}

	resp, err := client.Indices.Create(ctx, opensearchapi.IndicesCreateReq{
		Index: name,
		Body: bytes.NewReader(body),
	})
	if err != nil {
		return fmt.Errorf("failed to create indice %s: %w", name, err)
	}

	if !resp.Acknowledged {
		return fmt.Errorf("Indice request received but not acknowledged!")
	}

	log.Infof("Created indice '%s'", name)

	return nil
}

func hasIndice(ctx context.Context, name string) (bool, int, error) {
	resp, err := client.Indices.Exists(ctx, opensearchapi.IndicesExistsReq{
		Indices: []string{name},
	})
	if resp.StatusCode == http.StatusOK {
		return true, 0, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return false, 0, nil
	}

	if err != nil {
		return false, 0, err
	}

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	return false, 0, fmt.Errorf("Unexpected status code %d when checking indice %s: %s", resp.StatusCode, name, buf.Bytes())
}
