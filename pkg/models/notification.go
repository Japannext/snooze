package models

import (
	"fmt"
)

const NOTIFICATION_INDEX = "v2-notifications"

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

	ActiveCheckURL string `json:"activeCheckURL"`
}

func (item *Notification) GetID() string { return item.ID }
func (item *Notification) SetID(id string) { item.ID = id }

type Destination struct {
	// Name of the notification queue it will be sent to
	Queue string `json:"queue" yaml:"queue"`

	// Name of the profile (when notification have multiple profiles, e.g. destinations)
	Profile string `json:"profile" yaml:"profile"`
}

func (dest *Destination) String() string {
	return fmt.Sprintf("%s:%s", dest.Queue, dest.Profile)
}

func init() {
	index := IndexTemplate{
		Version: 0,
		IndexPatterns: []string{NOTIFICATION_INDEX},
		DataStream: map[string]map[string]string{"timestamp_field": {"name": "timestampMillis"}},
		Template: Indice{
			Settings: IndexSettings{1, 2},
			Mappings: IndexMapping{
				Properties: map[string]MappingProps{
					"timestampMillis": {Type: "date", Format: "epoch_millis"},
					"destination.kind": {Type: "keyword"},
					"destination.name": {Type: "keyword"},
				},
			},
		},
	}
	INDEXES = append(INDEXES, index)
}
