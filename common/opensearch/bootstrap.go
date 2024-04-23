package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	v2 "github.com/opensearch-project/opensearch-go/v2"
	api "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

type IndexSettings struct {
	NumberOfShards   int `json:"number_of_shards"`
	NumberOfReplicas int `json:"number_of_replicas"`
}

type MappingProps struct {
	Type   string                  `json:"type"`
	Fields map[string]MappingProps `json:"fields,omitempty"`
}

type IndexMapping struct {
	Properties map[string]MappingProps `json:"properties"`
}

type Indice struct {
	Name     string
	Settings IndexSettings  `json:"settings"`
	Mappings []IndexMapping `json:"mappings"`
}

type IndexTemplate struct {
	Name         string
	IndexPattern []string `json:"index_patterns"`
	Template     Indice   `json:"template"`
}

func indices() []IndexTemplate {
	numberOfShards := 3
	numberOfReplicas := 3
	settings := IndexSettings{numberOfShards, numberOfReplicas}
	return []IndexTemplate{
		{
			Name:         "snooze-log-v2",
			IndexPattern: []string{"snooze-log-v2-*"},
			Template: Indice{
				Settings: settings,
				Mappings: []IndexMapping{
					{
						Properties: map[string]MappingProps{
							"kind":       {Type: "keyword"},
							"timestamp":  {Type: "unsigned_long"},
							"group.hash": {Type: "byte"},
							"group.kv":   {Type: "object"},
							"resource":   {Type: "object"},
							"attributes": {Type: "object"},
							"body":       {Type: "object"},
							"snooze": {
								Type: "object",
								Fields: map[string]MappingProps{
									"snoozed": {Type: "boolean"},
									"name":    {Type: "keyword"},
								},
							},
						},
					},
				},
			},
		},
		{
			Name:         "snooze-group-v2",
			IndexPattern: []string{"snooze-group-v2-*"},
			Template: Indice{
				Settings: settings,
				Mappings: []IndexMapping{
					{
						Properties: map[string]MappingProps{
							"first_ts":  {Type: "unsigned_long"},
							"last_ts":   {Type: "unsigned_long"},
							"kv":        {Type: "object"},
							"last_body": {Type: "object"},
							"hits":      {Type: "integer"},
						},
					},
				},
			},
		},
	}
}

func (i *IndexTemplate) ensure(ctx context.Context, client *v2.Client) error {

	resp, err := api.IndicesExistsIndexTemplateRequest{Name: i.Name}.Do(ctx, client)
	if err != nil {
		return fmt.Errorf("Error while checking if index template '%s' exists: %w", i.Name, err)
	}
	if resp.StatusCode == 200 {
		return nil
	}
	if resp.StatusCode == 404 {
		body, err := json.Marshal(i)
		if err != nil {
			return err
		}
		_, err = api.IndicesPutIndexTemplateRequest{
			Name: i.Name,
			Body: bytes.NewReader(body),
		}.Do(ctx, client)
		if err != nil {
			return fmt.Errorf("Error while creating index template '%s': %w", i.Name, err)
		}
		return nil
	}
	return fmt.Errorf("Unexpected response status code when creating index %s: %d", i.Name, resp.StatusCode)
}

func (client *OpensearchClient) Bootstrap() error {
	ctx := context.Background()

	// Ensure all indexes exist
	for _, idx := range indices() {
		err := idx.ensure(ctx, client.Client)
		if err != nil {
			return err
		}
	}

	return nil
}
