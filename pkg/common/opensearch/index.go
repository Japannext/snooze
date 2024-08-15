package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

type IndexSettings struct {
	NumberOfShards   int `json:"number_of_shards"`
	NumberOfReplicas int `json:"number_of_replicas"`
}

type MappingProps struct {
	Type   string                  `json:"type,omitempty"`
	Format string				   `json:"format,omitempty"`
	Fields map[string]MappingProps `json:"fields,omitempty"`
}

type IndexMapping struct {
	Properties map[string]MappingProps `json:"properties"`
}

type Indice struct {
	Settings IndexSettings  `json:"settings"`
	Mappings IndexMapping `json:"mappings"`
}

type IndexTemplate struct {
	IndexPattern []string `json:"index_patterns"`
	DataStream DataStream `json:"data_stream"`
	Template     Indice   `json:"template"`
}

type DataStream struct {
	TimestampField TimestampField `json:"timestamp_field"`
}

type TimestampField struct {
	Name string `json:"name"`
}

func (lst *OpensearchLogStore) createIndex(ctx context.Context, name string, tpl IndexTemplate) error {
	body, err := json.Marshal(tpl)
	if err != nil {
		return fmt.Errorf("error marshaling index template: %s", err)
	}
	resp, err := lst.Client.IndexTemplate.Create(ctx, opensearchapi.IndexTemplateCreateReq{
		IndexTemplate: name,
		Body: bytes.NewReader(body),
	})
	if err != nil {
		return err
	}
	if !resp.Acknowledged {
		return fmt.Errorf("Index template request received but not acknowledged!")
	}
	return nil
}

func (lst *OpensearchLogStore) ensureIndex(ctx context.Context, name string, tpl IndexTemplate) {
	resp, err := lst.Client.IndexTemplate.Exists(ctx, opensearchapi.IndexTemplateExistsReq{
		IndexTemplate: name,
	})
	if resp.StatusCode == 200 {
		log.Infof("Index %s already exists", name)
		return
	}
	if resp.StatusCode == 404 {
		if err := lst.createIndex(ctx, name, tpl); err != nil {
			log.Fatalf("failed to create index template '%s': %s", name, err)
		}
		log.Infof("Successfully created index %s", name)
		return
	}
	if err != nil {
		log.Fatalf("Error while checking if index template '%s' exists: %s", name, err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		log.Fatal(err)
	}
	log.Fatalf("Unexpected status code %d when checking index %s: %s", resp.StatusCode, name, buf.Bytes())
}

func (lst *OpensearchLogStore) createDatastream(ctx context.Context, name string) error {
	resp, err := lst.Client.DataStream.Create(ctx, opensearchapi.DataStreamCreateReq{DataStream: name})
	if err != nil {
		return err
	}
	if !resp.Acknowledged {
		return fmt.Errorf("Datastream request received but not acknowledged!")
	}
	return nil
}

func (lst *OpensearchLogStore) ensureDatastream(ctx context.Context, name string) {
	ds, err := lst.Client.DataStream.Get(ctx, &opensearchapi.DataStreamGetReq{DataStreams: []string{name}})
	resp := ds.Inspect().Response
	if resp.StatusCode == 200 {
		log.Infof("Datastream '%s' already exists", name)
		return
	}
	if resp.StatusCode == 404 {
		if err := lst.createDatastream(ctx, name); err != nil {
			log.Fatalf("Failed to create datastream '%s': %s", name, err)
		}
		log.Infof("Successfully created datastream '%s'", name)
		return
	}
	if err != nil {
		log.Fatalf("Failed to fetch datastream '%s': %s", name, err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		log.Fatal(err)
	}
	log.Fatalf("Unexpected status code %d when checking datastream %s: %s", resp.StatusCode, name, buf.Bytes())
}
