package models

const ALERT_INDEX = "v2-alerts"

type Alert struct {
	Base

	Hash string `json:"hash"`

	Source Source `json:"source"`

	Identity map[string]string `json:"identity"`

	// Timestamps
	StartAt Time `json:"startAt"`
	EndAt Time `json:"endAt"`

	// Name of the alert
	AlertName string `json:"alertName"`
	// Group of the alert
	AlertGroup string `json:"alertGroup"`

	// Text representing the severity
	SeverityText string `json:"severityText,omitempty"`
	// Number representing the severity. Useful for filters (severity higher than a given value)
	SeverityNumber int32 `json:"severityNumber,omitempty"`

	// Text representing the severity
	//SeverityText string `json:"severityText,omitempty"`
	// Number representing the severity. Useful for filters (severity higher than a given value)
	//SeverityNumber int32 `json:"severityNumber,omitempty"`

	Description string `json:"description"`
	Summary string `json:"summary"`

	// Key-value representing the main resource identifiers.
	// Examples: host, pod, disk name, IP address
	Labels map[string]string `json:"labels,omitempty"`
}

type AlertStatus struct {
	ID string `json:"id"`
	LastCheck Time `json:"lastCheck"`
	NextCheck Time `json:"nextCheck"`
}

type AlertUpdate struct {
	Document *Alert `json:"doc"`
	Upsert *Alert `json:"upsert"`
}

func init() {
	index := IndexTemplate{
		Version: 0,
		IndexPatterns: []string{ALERT_INDEX},
		DataStream: map[string]map[string]string{"timestamp_field": {"name": "startsAt"}},
		Template: Indice{
			Settings: IndexSettings{1, 2},
			Mappings: IndexMapping{
				Properties: map[string]MappingProps{
					"hash": {Type: "keyword"},
					"startsAt": {Type: "date", Format: "epoch_millis"},
					"endsAt": {Type: "date", Format: "epoch_millis"},
					"source.kind": {Type: "keyword"},
					"source.name": {Type: "keyword"},
					"labels": {Type: "object"},
				},
			},
		},
	}
	INDEXES = append(INDEXES, index)
}
