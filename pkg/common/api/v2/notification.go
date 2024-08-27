package v2

import (
	"fmt"
)

type NotificationResults struct {
	Items []Notification `json:"items"`
	Total int `json:"total"`
}

type Notification struct {
	ID string `json:"id,omitempty"`
	TimestampMillis   uint64 `json:"timestampMillis"`
	Destination Destination `json:"destination"`
	Acknowledged bool `json:"acknowledged"`
	AlertUID string `json:"alertUID,omitempty"`
	LogUID string `json:"logUID,omitempty"`
	Body map[string]string `json:"body"`

	Labels map[string]string `json:"labels"`
	DocumentationURL string `json:"documentationURL"`
	SnoozeURL string `json:"snoozeURL"`
}

type Destination struct {
	// Name of the notification queue it will be sent to
	Queue string `json:"queue" yaml:"queue"`

	// Name of the profile (when notification have multiple profiles, e.g. destinations)
	Profile string `json:"profile" yaml:"profile"`
}

func (dest *Destination) String() string {
	return fmt.Sprintf("%s:%s", dest.Queue, dest.Profile)
}
