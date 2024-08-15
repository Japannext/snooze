package v2

import (
	"time"
)

type TimeRange struct {
	Start time.Time `json:"start"`
	End time.Time `json:"end"`
}
