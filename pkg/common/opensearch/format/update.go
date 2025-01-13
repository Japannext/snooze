package format

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const updateAction Action = "update"

type Update struct {
	Index       string
	ID          string
	Doc         interface{}
	Upsert      interface{}
	DocAsUpsert bool
}

type updateWrapper struct {
	Doc         json.RawMessage  `json:"doc,omitempty"`
	DocAsUpsert bool             `json:"doc_as_upsert,omitempty"`
	Upsert      *json.RawMessage `json:"upsert,omitempty"`
}

func (a *Update) Serialize() ([]byte, error) {
	var buf bytes.Buffer

    meta := Metadata{Index: a.Index, ID: a.ID}

    header, err := MarshalHeader(updateAction, meta)
    if err != nil {
        return []byte{}, fmt.Errorf("error serializing header: %w", err)
    }

	buf.Write(header)
	buf.WriteString("\n")

	w := &updateWrapper{
		DocAsUpsert: a.DocAsUpsert,
	}
	if a.Doc != nil {
		w.Doc, err = json.Marshal(a.Doc)
		if err != nil {
			return []byte{}, fmt.Errorf("error serializing document `%+v`: %w", a.Doc, err)
		}
	}

	if a.Upsert != nil {
		upsert, err := json.Marshal(a.Upsert)
		if err != nil {
			return []byte{}, fmt.Errorf("error serializing upsert `%+v`: %w", a.Upsert, err)
		}

		upsertData := json.RawMessage(upsert)
		w.Upsert = &upsertData
	}

	data, err := json.Marshal(w)
	if err != nil {
		return []byte{}, fmt.Errorf("error serializing body `%+v`: %s", w, err)
	}

	buf.Write(data)
	buf.WriteString("\n")

	return buf.Bytes(), nil
}
