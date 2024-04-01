package opensearch

import (
  "bytes"
  "context"
  "encoding/json"

  api "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
  dsl "github.com/mottaquikarim/esquerydsl"

  "github.com/japannext/snooze/common/types"
)

func (db *Database) LogV2Insert(ctx context.Context, l types.LogV2) error {
  b, err := json.Marshal(l)
  if err != nil {
    return err
  }
  _, err = api.IndexRequest{Index: "snooze-log-v2", Body: bytes.NewReader(b)}.Do(ctx, db.Client)
  if err != nil {
    return err
  }
  return nil
}

func (db *Database) LogV2Search(ctx context.Context, search, sortBy string, pageNb, perPage int) ([]types.LogV2, error) {
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
    return []types.LogV2{}, err
  }
  resp, err := api.SearchRequest{
    Index: []string{"snooze-log-v2"},
    Body: bytes.NewReader(body),
  }.Do(ctx, db.Client)
  if err != nil {
    return []types.LogV2{}, err
  }
  var buf bytes.Buffer
  buf.ReadFrom(resp.Body)
  var ll []types.LogV2
  json.Unmarshal(buf.Bytes(), &ll)
  return ll, nil
}
