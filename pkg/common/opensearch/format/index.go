package format

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const indexAction Action = "index"

type Index struct {
	Index string
	ID    string
	Item  interface{}
}

func (a *Index) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	meta := Metadata{Index: a.Index}
	if a.ID != "" {
		meta.ID = a.ID
	}

	data, err := MarshalHeader(indexAction, meta)
	if err != nil {
		return []byte{}, fmt.Errorf("error serializing header: %w", err)
	}

	buf.Write(data)
	buf.WriteString("\n")

	body, err := json.Marshal(a.Item)
	if err != nil {
		return []byte{}, fmt.Errorf("error serializing body `%+v`: %w", a.Item, err)
	}

	buf.Write(body)
	buf.WriteString("\n")

	return buf.Bytes(), nil
}
