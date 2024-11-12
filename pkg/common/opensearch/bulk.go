package opensearch

import (
	"context"
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

func Bulk(ctx context.Context, req opensearchapi.BulkReq) (*opensearchapi.BulkResp, error) {
       return client.Bulk(ctx, req)
}

/*
func bulkRespToError(resp *opensearchapi.BulkResp) error {
       var buf strings.Builder
       for i, v := range resp.Items {
               for k, r := range v {
                       if r.Error != nil {
                               msg := fmt.Sprintf("[#%d:%s] type='%s' reason='%s'\n", i, k, r.Error.Type, r.Error.Reason)
                               buf.WriteString(msg)
                       }
               }
       }
       return fmt.Errorf("error in bulk log:\n%s", buf.String())
}
*/

type Action = string

const (
	INDEX_ACTION  Action = "index"
	UPDATE_ACTION Action = "update"
	CREATE_ACTION Action = "create"
	DELETE_ACTION Action = "delete"
)

type BulkHeader = map[Action]Metadata

type Metadata struct {
	ID string `json:"_id,omitempty"`
	Index string `json:"_index,omitempty"`
}

// CREATE action

type CreateAction struct {
	Index string
	Item interface{}
}

func (a *CreateAction) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	header := BulkHeader(map[Action]Metadata{
		CREATE_ACTION: {Index: a.Index},
	})
	data, err := json.Marshal(header)
	if err != nil {
		return []byte{}, fmt.Errorf("error marshalling `%+v`: %s", header, err)
	}
	buf.Write(data)
	buf.WriteString("\n")
	body, err := json.Marshal(a.Item)
	if err != nil {
		return []byte{}, fmt.Errorf("error marshalling `%+v`: %s", a.Item, err)
	}
	buf.Write(body)
	buf.WriteString("\n")
	return buf.Bytes(), nil
}

func Create(index string, item interface{}) *CreateAction {
	return &CreateAction{
		Index: index,
		Item: item,
	}
}

type IndexAction struct {
	Index string
	Item interface{}
}
func (a *IndexAction) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	header := BulkHeader(map[Action]Metadata{
		INDEX_ACTION: {Index: a.Index},
	})
	data, err := json.Marshal(header)
	if err != nil {
		return []byte{}, fmt.Errorf("error marshalling `%+v`: %s", header, err)
	}
	buf.Write(data)
	buf.WriteString("\n")
	body, err := json.Marshal(a.Item)
	if err != nil {
		return []byte{}, fmt.Errorf("error marshalling `%+v`: %s", a.Item, err)
	}
	buf.Write(body)
	buf.WriteString("\n")
	return buf.Bytes(), nil
}

// UPDATE action

type UpdateAction struct {
	Index string
	ID string
	Doc interface{}
	Upsert interface{}
	DocAsUpsert bool
}

type updateWrapper struct {
	Doc json.RawMessage `json:"doc,omitempty"`
	DocAsUpsert bool `json:"doc_as_upsert,omitempty"`
	Upsert *json.RawMessage `json:"upsert,omitempty"`
}

func (a *UpdateAction) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	header := BulkHeader(map[Action]Metadata{
		UPDATE_ACTION: {Index: a.Index, ID: a.ID},
	})
	headerData, err := json.Marshal(header)
	if err != nil {
		return []byte{}, fmt.Errorf("error marshalling `%+v`: %s", header, err)
	}
	buf.Write(headerData)
	buf.WriteString("\n")


	w := &updateWrapper{
		DocAsUpsert: a.DocAsUpsert,
	}
	if a.Doc != nil {
		w.Doc, err = json.Marshal(a.Doc)
		if err != nil {
			return []byte{}, fmt.Errorf("error marhsalling `%+v`: %s", a.Doc, err)
		}
	}

	if a.Upsert != nil {
		upsert, err := json.Marshal(a.Upsert)
		if err != nil {
			return []byte{}, fmt.Errorf("error marhsalling `%+v`: %s", a.Upsert, err)
		}
		w.Upsert = pointer(json.RawMessage(upsert))
	}

	data, err := json.Marshal(w)
	if err != nil {
		return []byte{}, fmt.Errorf("error marhsalling `%+v`: %s", w, err)
	}
	buf.Write(data)
	buf.WriteString("\n")

	return buf.Bytes(), nil
}

// DELETE action
type DeleteAction struct {
	Index string
	ID string
}

func (a *DeleteAction) Serialize() ([]byte, error) {
	header := BulkHeader(map[Action]Metadata{
		DELETE_ACTION: {Index: a.Index, ID: a.ID},
	})
	data, err := json.Marshal(header)
	if err != nil {
		return []byte{}, fmt.Errorf("error marshalling `%+v`: %s", header, err)
	}
	return data, nil
}

func Delete(index, id string) *DeleteAction {
	return &DeleteAction{
		Index: index,
		ID: id,
	}
}
