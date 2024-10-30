package models

import (
	"fmt"
)

const NOTIFICATION_INDEX = "v2-notifications"

type Notification struct {
	ID string `json:"id,omitempty"`
	Timestamp Timestamp `json:"timestamp"`
	Destination Destination `json:"destination"`
	Acknowledged bool `json:"acknowledged"`

	// The type of notification. Supported: "log", "alert"
	Type string `json:"type"`
	// The UID of the log of alert referenced by the notification
	ItemID string `json:"itemID,omitempty"`
	// The identity of the log/alert
	Identity map[string]string `json:"identity,omitempty"`
	// The message contained in the log or the summary of the alert
	Message string `json:"message,omitempty"`

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
		DataStream: map[string]map[string]string{"timestamp_field": {"name": "timestamp.display"}},
		Template: Indice{
			Settings: IndexSettings{1, 2},
			Mappings: IndexMapping{
				Properties: map[string]MappingProps{
					"timestamp.display": {Type: "date", Format: "epoch_millis"},
					"timestamp.actual": {Type: "date", Format: "epoch_millis"},
					"timestamp.observed": {Type: "date", Format: "epoch_millis"},
					"timestamp.processed": {Type: "date", Format: "epoch_millis"},
					"timestamp.warning": {Type: "keyword"},
					"destination.kind": {Type: "keyword"},
					"destination.name": {Type: "keyword"},
				},
			},
		},
	}
	INDEXES = append(INDEXES, index)
}
