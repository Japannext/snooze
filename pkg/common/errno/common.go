package errno

import (
	"fmt"
)

func JSONMarshalError(item interface{}, err error) error {
	return fmt.Errorf("failed to marshal object type %T (%+v): %w", item, item, err)
}

func JSONUnmarshalError(data []byte, err error) error {
	return fmt.Errorf("failed to unmarshal data `%s`: %w", data[:20], err)
}
