package opensearch

import (
  "bufio"
  "bytes"
  "context"
  "encoding/json"
  "fmt"
  "io"

  api "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
  dsl "github.com/mottaquikarim/esquerydsl"
  // util "github.com/opensearch-project/opensearch-go/opensearchutil"

  "github.com/japannext/snooze/common/api/v2"
)

const (
  ALERT_EVENTS_V2 = "alert-events-v2"
)

func (client *OpensearchClient) SearchAlertEvent(ctx context.Context, search, sortBy string, pageNb, perPage int) ([]v2.Alert, error) {
  query := dsl.QueryDoc{
    From: pageNb * perPage,
    Size: perPage,
    PageSize: perPage,
    Sort: sorts(byTimestamp),
    Or: []dsl.QueryItem{
      {Field: "body", Value: search, Type: dsl.Match},
    },
  }
  body, err := json.Marshal(query)
  if err != nil {
    return []v2.Alert{}, err
  }
  req := api.SearchRequest{
    Index: []string{ALERT_EVENTS_V2},
    Body: bytes.NewReader(body),
  }
  resp, err := req.Do(ctx, client.Client)
  if err != nil {
    return []v2.Alert{}, err
  }
  var buf bytes.Buffer
  if _, err := buf.ReadFrom(resp.Body); err != nil {
    return []v2.Alert{}, err
  }
  var items []v2.Alert
  if err := json.Unmarshal(buf.Bytes(), &items); err != nil {
    return []v2.Alert{}, err
  }
  return items, nil
}

func (client *OpensearchClient) InsertAlertEvent(ctx context.Context, a *v2.Alert) error {
  b, err := json.Marshal(a)
  if err != nil {
    return err
  }
  body := bytes.NewReader(b)
  req := api.IndexRequest{
    Index: ALERT_EVENTS_V2,
    Body: body,
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
  client *OpensearchClient
  index string
  payload *bufio.ReadWriter
}

func NewBulkIndexer(client *OpensearchClient, index string) *BulkIndexer {
  var b bytes.Buffer
  r := bufio.NewReader(&b)
  w := bufio.NewWriter(&b)
  rw := bufio.NewReadWriter(r, w)
  return &BulkIndexer{
    client: client,
    index: index,
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

func (client *OpensearchClient) BulkInsertAlertEvent(ctx context.Context, items []v2.Alert) error {
  bulk := NewBulkIndexer(client, "alert-events-v2")
  for _, item := range items {
    bulk.Add(item)
  }
  if err := bulk.Flush(ctx); err != nil {
    return err
  }
  return nil
}
