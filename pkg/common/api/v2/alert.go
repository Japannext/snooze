package v2

type Alert struct {
	ID string `json:"id,omitempty"`

	Hash string `json:"hash"`

	Source Source `json:"source"`

	// Timestamps
	StartsAt uint64 `json:"startsAt"`
	EndsAt uint64 `json:"endsAt"`

	// Text representing the severity
	//SeverityText string `json:"severityText,omitempty"`
	// Number representing the severity. Useful for filters (severity higher than a given value)
	//SeverityNumber int32 `json:"severityNumber,omitempty"`

	Message string `json:"message"`
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

type AlertResults struct {
	Items []Alert `json:"items"`
	Total int `json:"total"`
	More bool `json:"more"`
}
