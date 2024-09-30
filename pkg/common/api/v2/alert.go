package v2

const ALERT_INDEX = "v2-alerts"

type Alert struct {
	ID string `json:"id,omitempty"`

	Hash string `json:"hash"`

	Source Source `json:"source"`

	Identity map[string]string `json:"identity"`

	// Timestamps
	StartsAt uint64 `json:"startsAt"`
	EndsAt uint64 `json:"endsAt"`

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

	// Mute the alert. This may skip notifications, or skip even display
	Mute Mute `json:"mute"`
}

func (item *Alert) GetID() string { return item.ID }
func (item *Alert) SetID(id string) { item.ID = id }

type AlertStatus struct {
	ID string `json:"id"`
	LastCheck uint64 `json:"lastCheck"`
	NextCheck uint64 `json:"nextCheck"`
}

type AlertUpdate struct {
	Document *Alert `json:"doc"`
	Upsert *Alert `json:"upsert"`
}

func init() {
	index := IndexTemplate{
		Version: 0,
		IndexPatterns: []string{ALERT_INDEX},
		DataStream: map[string]map[string]string{"timestamp_field": {"name": "timestampMillis"}},
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
