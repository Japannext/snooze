package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type timeTest struct {
	Name string
	Data []byte
	Time Time
}

func timeParse(txt string) Time {
	ts, err := time.Parse(time.RFC3339, txt)
	if err != nil {
		panic(err)
	}
	return Time{Time: ts}
}

const FMT = "2006-01-02T15:04:05.999+00:00"

var marshalTimes = []timeTest{
	{Name: "unix zero", Data: []byte("0"), Time: timeParse("1970-01-01T00:00:00+00:00")},
	{Name: "today", Data: []byte("1730251555000"), Time: timeParse("2024-10-30T10:25:55+09:00")},
}

func TestTimeMarshal(t *testing.T) {
	for _, tt := range marshalTimes {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := json.Marshal(tt.Time)
			if err != nil {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.Data, data)
		})
	}
}

func TestTimeUnmarshal(t *testing.T) {
	for _, tt := range marshalTimes {
		t.Run(tt.Name, func(t *testing.T) {
			ts := Time{}
			err := json.Unmarshal(tt.Data, &ts)
			if err != nil {
				assert.NoError(t, err)
			}
			assert.WithinDuration(t, tt.Time.Time, ts.Time, 0)
		})
	}
}
