package format

import (
	"encoding/json"
	"fmt"
)

type Delete struct {
	Index string
	ID    string
}

func (a *Delete) Serialize() ([]byte, error) {
	header := BulkHeader(map[Action]Metadata{
		DELETE_ACTION: {Index: a.Index, ID: a.ID},
	})
	data, err := json.Marshal(header)
	if err != nil {
		return []byte{}, fmt.Errorf("error marshalling `%+v`: %s", header, err)
	}
	return data, nil
}
