package opensearch

import (
	"bytes"
	"context"
	"fmt"

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
	if resp.StatusCode == 200 {
		log.Infof("Datastream '%s' already exists", name)
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
	return false, fmt.Errorf("Unexpected status code %d when checking datastream %s: %s", resp.StatusCode, name, buf.Bytes())
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
