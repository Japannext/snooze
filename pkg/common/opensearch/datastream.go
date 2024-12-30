package opensearch

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

func createDatastream(ctx context.Context, name string) error {
	resp, err := client.DataStream.Create(ctx, opensearchapi.DataStreamCreateReq{DataStream: name})
	if err != nil {
		return err
	}
	if !resp.Acknowledged {
		return fmt.Errorf("Datastream request received but not acknowledged!")
	}
	log.Infof("Created datastream %s", name)
	return nil
}

func hasDatastream(ctx context.Context, name string) (bool, error) {
	ds, err := client.DataStream.Get(ctx, &opensearchapi.DataStreamGetReq{DataStreams: []string{name}})
	resp := ds.Inspect().Response
	if resp.StatusCode == http.StatusOK {
		log.Infof("Datastream '%s' already exists", name)
		return true, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return false, fmt.Errorf("cannot ready body: %w", err)
	}

	return false, fmt.Errorf("status code %d in datastream %s: %s", resp.StatusCode, name, buf.Bytes())
}

func ensureDatastream(ctx context.Context, name string) {
	found, err := hasDatastream(ctx, name)
	if err != nil {
		log.Fatal(err)
	}
	if found {
		return
	}
	if err := createDatastream(ctx, name); err != nil {
		log.Fatal(err)
	}
}
