package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"github.com/japannext/snooze/pkg/models"
)

func createIndex(ctx context.Context, name string, tpl models.IndexTemplate) error {
	body, err := json.Marshal(tpl)
	if err != nil {
		return fmt.Errorf("error marshaling index template: %s", err)
	}
	log.Debugf("Template: %s", body)
	resp, err := client.IndexTemplate.Create(ctx, opensearchapi.IndexTemplateCreateReq{
		IndexTemplate: name,
		Body:          bytes.NewReader(body),
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

func hasIndex(ctx context.Context, name string) (bool, int, error) {
	resp, err := client.IndexTemplate.Exists(ctx, opensearchapi.IndexTemplateExistsReq{
		IndexTemplate: name,
	})
	if resp.StatusCode == 200 {
		r, err := client.IndexTemplate.Get(ctx, &opensearchapi.IndexTemplateGetReq{
			IndexTemplates: []string{name},
		})
		if err != nil {
			return true, 0, err
		}

		details := r.IndexTemplates[0]
		version := details.IndexTemplate.Version
		return true, version, nil
	}
	if resp.StatusCode == 404 {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}
	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	return false, 0, fmt.Errorf("Unexpected status code %d when checking index %s: %s", resp.StatusCode, name, buf.Bytes())
}

func ensureIndex(ctx context.Context, name string, tpl models.IndexTemplate) {
	found, version, err := hasIndex(ctx, name)
	if err != nil {
		log.Fatal(err)
	}
	if found && tpl.Version <= version {
		log.Debugf("index '%s' (version=%d) already exist", name, version)
		return
	}
	if err := createIndex(ctx, name, tpl); err != nil {
		log.Fatal(err)
	}
}
