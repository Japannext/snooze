package format

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const createAction Action = "create"

type Create struct {
	Index string
	Item  interface{}
}

func (a *Create) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	meta := Metadata{Index: a.Index}
	data, err := MarshalHeader(createAction, meta)
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
