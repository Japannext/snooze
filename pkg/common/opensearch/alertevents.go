package opensearch

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	dsl "github.com/mottaquikarim/esquerydsl"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	// util "github.com/opensearch-project/opensearch-go/opensearchutil"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

const (
	ALERT_EVENTS_V2 = "alert-events-v2"
)

func (client *OpensearchLogStore) Search(ctx context.Context, query string, pagination *api.Pagination) ([]api.Alert, error) {
	dslQuery := dsl.QueryDoc{
		From:     pagination.PageNumber * pagination.PerPage,
		Size:     pagination.PerPage,
		PageSize: pagination.PerPage,
		Sort:     sorts(byTimestamp),
		Or: []dsl.QueryItem{
			{Field: "body", Value: query, Type: dsl.Match},
		},
	}
	body, err := json.Marshal(dslQuery)
	if err != nil {
		return []api.Alert{}, err
	}
	req := opensearchapi.SearchRequest{
		Index: []string{ALERT_EVENTS_V2},
		Body:  bytes.NewReader(body),
	}
	resp, err := req.Do(ctx, client.Client)
	if err != nil {
		return []api.Alert{}, err
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return []api.Alert{}, err
	}
	var items []api.Alert
	if err := json.Unmarshal(buf.Bytes(), &items); err != nil {
		return []api.Alert{}, err
	}
	return items, nil
}

func (client *OpensearchLogStore) Store(alert *api.Alert) error {
	ctx := context.Background()
	b, err := json.Marshal(alert)
	if err != nil {
		return err
	}
	body := bytes.NewReader(b)
	req := opensearchapi.IndexRequest{
		Index: ALERT_EVENTS_V2,
		Body:  body,
	}
	if _, err := req.Do(ctx, client.Client); err != nil {
		return err
	}
	return nil
}

type PartialError struct {
	failed uint32
	errors []error
}

func (p *PartialError) Add(err error) {
	p.failed++
	p.errors = append(p.errors, err)
}

func (p *PartialError) Empty() bool {
	return (p.failed == 0)
}

func (p *PartialError) Error() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("Partial error for %d messages:", p.failed))
	for _, err := range p.errors {
		b.WriteString(err.Error())
	}
	return b.String()
}

type Response struct {
	Took   int  `json:"took"`
	Errors bool `json:"errors"`
	Items  []struct {
		Delete struct {
			Index   string `json:"_index"`
			Id      string `json:"_id"`
			Version int    `json:"_version"`
			Result  string `json:"result"`
			Shards  struct {
				Total      int `json:"total"`
				Successful int `json:"successful"`
				Failed     int `json:"failed"`
			} `json:"_shards"`
			SeqNo       int `json:"_seq_no"`
			PrimaryTerm int `json:"_primary_term"`
			Status      int `json:"status"`
		} `json:"delete,omitempty"`
	} `json:"items"`
}

type Operation struct {
	Name string
	Body []byte
}

type BulkIndexer struct {
	client  *OpensearchLogStore
	index   string
	payload *bufio.ReadWriter
}

func NewBulkIndexer(client *OpensearchLogStore, index string) *BulkIndexer {
	var b bytes.Buffer
	r := bufio.NewReader(&b)
	w := bufio.NewWriter(&b)
	rw := bufio.NewReadWriter(r, w)
	return &BulkIndexer{
		client:  client,
		index:   index,
		payload: rw,
	}
}

func (bulk *BulkIndexer) Add(obj any) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	bulk.payload.WriteString(fmt.Sprintf(`{"index": {"_index": %s}}`, bulk.index))
	bulk.payload.Write([]byte("\n"))
	bulk.payload.Write(b)
	bulk.payload.Write([]byte("\n"))
	return nil
}

func (bulk *BulkIndexer) Flush(ctx context.Context) error {
	// perr := &PartialError{}
	res, err := bulk.client.Client.Bulk(bulk.payload, bulk.client.Client.Bulk.WithContext(ctx))
	if err != nil {
		return err
	}
	rb, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response Response
	if err := json.Unmarshal(rb, &response); err != nil {
		return err
	}
	return nil
}

func (client *OpensearchLogStore) BulkInsertAlertEvent(ctx context.Context, items []api.Alert) error {
	bulk := NewBulkIndexer(client, "alert-events-v2")
	for _, item := range items {
		bulk.Add(item)
	}
	if err := bulk.Flush(ctx); err != nil {
		return err
	}
	return nil
}
