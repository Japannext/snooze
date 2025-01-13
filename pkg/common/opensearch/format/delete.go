package format

import (
	"fmt"
)

const deleteAction Action = "delete"

type Delete struct {
	Index string
	ID    string
}

func (a *Delete) Serialize() ([]byte, error) {
    meta := Metadata{Index: a.Index, ID: a.ID}

    data, err := MarshalHeader(deleteAction, meta)
    if err != nil {
        return []byte{}, fmt.Errorf("error serializing header: %w", err)
    }

	return data, nil
}
