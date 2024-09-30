package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func createIndex(ctx context.Context, name string, tpl api.IndexTemplate) error {
	body, err := json.Marshal(tpl)
	if err != nil {
		return fmt.Errorf("error marshaling index template: %s", err)
	}
	log.Debugf("Template: %s", body)
	resp, err := client.IndexTemplate.Create(ctx, opensearchapi.IndexTemplateCreateReq{
		IndexTemplate: name,
		Body: bytes.NewReader(body),
	})
	if err != nil {
		return err
	}
	if !resp.Acknowledged {
		return fmt.Errorf("Index template request received but not acknowledged!")
	}
	log.Infof("Created index %s", name)
	return nil
}

func hasIndex(ctx context.Context, name string) (bool, error) {
	resp, err := client.IndexTemplate.Exists(ctx, opensearchapi.IndexTemplateExistsReq{
		IndexTemplate: name,
	})
	if resp.StatusCode == 200 {
		return true, nil
	}
	if resp.StatusCode == 404 {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	return false, fmt.Errorf("Unexpected status code %d when checking index %s: %s", resp.StatusCode, name, buf.Bytes())
}

func ensureIndex(ctx context.Context, name string, tpl api.IndexTemplate) {
	found, err := hasIndex(ctx, name)
	if err != nil {
		log.Fatal(err)
	}
	if found {
		return
	}
	if err := createIndex(ctx, name, tpl); err != nil {
		log.Fatal(err)
	}
}
