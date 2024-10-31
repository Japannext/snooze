package models

import (
	"fmt"
	"strconv"
	"time"

	// Important for making timezone work as expected
	_ "time/tzdata"
)

// Wrapper for marshal/unmarshal all times in the same way
type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(t.UnixMilli()))), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	i, err := strconv.Atoi(string(data))
	if err != nil {
		return fmt.Errorf("failed to unmarshal time `%s`: %s", data, err)
	}
	tt := time.UnixMilli(int64(i))
	t.Time = tt.Local()
	return nil
}

func TimeNow() Time {
	return Time{time.Now()}
}
