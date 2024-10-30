package models

type Timestamp struct {
	// The timestamp that will be displayed in the webUI
	Display Time `json:"display,omitempty"`
	// Actual time of the event when supported
	Actual Time `json:"actual,omitempty"`
	// When the event was observed by the earliest relay
	// supporting it (opentelemetry or snooze)
	Observed Time `json:"observed,omitempty"`
	// Time when the log was processed (written in the database)
	Processed Time `json:"processed,omitempty"`
	// Some warning when there are timestamp inconsistencies
	Warning string `json:"warning,omitempty"`
}

