package format

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Create struct {
	Index string
	Item interface{}
}

func (a *Create) Serialize() ([]byte, error) {
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
