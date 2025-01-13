package format

import (
	"encoding/json"
	"fmt"
)

type Action = string

type Metadata struct {
	ID    string `json:"_id,omitempty"`
	Index string `json:"_index,omitempty"`
}

func MarshalHeader(action Action, meta Metadata) ([]byte, error) {
	header := map[Action]Metadata{action: meta}

	data, err := json.Marshal(header)
	if err != nil {
		return []byte{}, fmt.Errorf("error marshalling `%+v`: %w", header, err)
	}

	return data, nil
}
