package opensearch

import (
	// "bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	// "io"
	// "time"

	dsl "github.com/mottaquikarim/esquerydsl"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	// "github.com/opensearch-project/opensearch-go/v4/opensearchutil"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

type OpensearchLogStore struct {
	Client *opensearchapi.Client
}

func NewLogStore() *OpensearchLogStore {
	cfg, err := initConfig()
	if err != nil {
		log.Fatal(err)
	}

	client, err := opensearchapi.NewClient(opensearchapi.Config{Client: cfg})
	if err != nil {
		log.Fatal(err)
	}

	return &OpensearchLogStore{client}
}

func (lst *OpensearchLogStore) GetLog(uid string) (*api.Log, error) {
	ctx := context.Background()
	dslQuery := dsl.QueryDoc{
		PageSize: 1,
		Or: []dsl.QueryItem{
			{Field: "_id", Value: uid, Type: dsl.Term},
		},
	}
	body, err := json.Marshal(dslQuery)
	if err != nil {
		return nil, err
	}
	req := &opensearchapi.SearchReq{
		Indices: []string{"log-v2"},
		Body:  bytes.NewReader(body),
	}
	resp, err := lst.Client.Search(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Errors {
		return nil, fmt.Errorf("opensearch returned an error: %s", "")
	}
	if resp.Hits.Total.Value > 0 {
		hit := resp.Hits.Hits[0]
		var item *api.Log
		if err := json.Unmarshal(hit.Fields, &item); err != nil {
			return nil, fmt.Errorf("cannot unmarshal to Log: %w", err)
		}
		item.ID = hit.ID
		return item, nil
	}
	return nil, nil
}

type Range struct {
	Gte *uint64 `json:"gte,omitempty"`
	Lte *uint64 `json:"lte,omitempty"`
}

func addTimeRange(doc *dsl.QueryDoc, timerange api.TimeRange) {
	var r = Range{}
	if timerange.Start > 0 {
		r.Gte = &timerange.Start
	}
	if timerange.End > 0 {
		r.Lte = &timerange.End
	}
	if timerange.Start > 0 || timerange.End > 0 {
		item := dsl.QueryItem{Type: dsl.Range, Field: "timestampMillis", Value: r}
		doc.And = append(doc.And, item)
	}
}

func addPagination(doc *dsl.QueryDoc, params *opensearchapi.SearchParams, pagination api.Pagination) {
	params.From = pointer((pagination.Page - 1) * pagination.Size)
	params.Size = pointer(pagination.Size)
	sort := map[string]string{}
	if pagination.Ascending {
		sort[pagination.OrderBy] = "asc"
	} else {
		sort[pagination.OrderBy] = "desc"
	}
	doc.Sort = []map[string]string{sort}
}

func addSearch(doc *dsl.QueryDoc, search string) {
	if search != "" {
		item := dsl.QueryItem{Type: dsl.Match, Field: "message", Value: search}
		doc.And = append(doc.And, item)
	}
}

func (lst *OpensearchLogStore) SearchLogs(ctx context.Context, search string, timerange api.TimeRange, pagination api.Pagination) (api.LogResults, error) {
	var res api.LogResults
	var doc = &dsl.QueryDoc{}
	var params = &opensearchapi.SearchParams{}

	addTimeRange(doc, timerange)
	addPagination(doc, params, pagination)
	addSearch(doc, search)

	body, err := json.Marshal(doc)
	if err != nil {
		return res, fmt.Errorf("invalid request body (%+v): %w", doc, err)
	}
	req := &opensearchapi.SearchReq{
		Indices: []string{"log-v2"},
		Params: *params,
		Body:  bytes.NewReader(body),
	}
	resp, err := lst.Client.Search(ctx, req)
	if err != nil {
		return res, err
	}
	if resp.Errors {
		return res, fmt.Errorf("opensearch returned an error: %s", "")
	}

	res.Items = make([]api.Log, len(resp.Hits.Hits))
	for i, hit := range resp.Hits.Hits {
		if err := json.Unmarshal(hit.Source, &res.Items[i]); err != nil {
			return res, fmt.Errorf("failed to unmarshal message %s: %w", hit.Source, err)
		}
		res.Items[i].ID = hit.ID
	}
	res.Total = resp.Hits.Total.Value
	if resp.Hits.Total.Relation != "eq" {
		res.More = true
	}
	return res, nil
}

func (lst *OpensearchLogStore) SearchNotifications(ctx context.Context, search string, timerange api.TimeRange, pagination api.Pagination) (api.NotificationResults, error) {
	var res api.NotificationResults
	resp, err := lst.search(ctx, "notification-v2", search, timerange, pagination)
	if err != nil {
		return res, err
	}
	res.Items = make([]api.Notification, len(resp.Hits.Hits))
	for i, hit := range resp.Hits.Hits {
		if err := json.Unmarshal(hit.Source, &res.Items[i]); err != nil {
			return res, fmt.Errorf("failed to unmarshal message %s: %w", hit.Source, err)
		}
		res.Items[i].ID = hit.ID
	}
	res.Total = resp.Hits.Total.Value
	return res, nil
}

func (lst OpensearchLogStore) search(ctx context.Context, index, search string, timerange api.TimeRange, pagination api.Pagination) (*opensearchapi.SearchResp, error) {
	var doc = &dsl.QueryDoc{}
	var params = &opensearchapi.SearchParams{}

	addTimeRange(doc, timerange)
	addPagination(doc, params, pagination)
	addSearch(doc, search)

	body, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("invalid request body (%+v): %w", doc, err)
	}
	req := &opensearchapi.SearchReq{
		Indices: []string{index},
		Params: *params,
		Body:  bytes.NewReader(body),
	}
	resp, err := lst.Client.Search(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Errors {
		return nil, fmt.Errorf("opensearch returned an error: %s", "")
	}
	return resp, nil
}

func (lst *OpensearchLogStore) store(index string, item interface{}) error {
	ctx := context.Background()
	b, err := json.Marshal(item)
	if err != nil {
		return err
	}
	body := bytes.NewReader(b)
	req := opensearchapi.IndexReq{
		Index: index,
		Body:  body,
	}
	resp, err := lst.Client.Index(ctx, req)
	if err != nil {
		return err
	}
	if resp.Shards.Successful == 0 {
		return fmt.Errorf("Failed to index object to '%s': %s", index, resp.Result)
	}
	log.Debugf("inserted into `%s`: %s", index, b)
	return nil
}

func (lst *OpensearchLogStore) StoreLog(item *api.Log) error {
	return lst.store("log-v2", item)
}

func (lst *OpensearchLogStore) StoreNotification(item *api.Notification) error {
	return lst.store("notification-v2", item)
}

/*

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

func (lst *OpensearchLogStore) BulkInsertLogEvent(ctx context.Context, items []api.Log) error {
	bulk := NewBulkIndexer(lst, "alert-events-v2")
	for _, item := range items {
		bulk.Add(item)
	}
	if err := bulk.Flush(ctx); err != nil {
		return err
	}
	return nil
}
*/
