package format

import (
	"bytes"
	"encoding/json"
	"fmt"
)

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
	header := BulkHeader(map[Action]Metadata{
		INDEX_ACTION: meta,
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
