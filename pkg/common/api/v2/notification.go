package v2

type Destination struct {
	// Type of destination (e.g mail/googlechat/patlite/browser)
	Kind string `json:"kind"`
	// Name of the instance of destination backend (e.g. prod/example.com)
	Name string `json:"name"`
}

type NotificationResults struct {
	Items []Notification
	Total int
}

type Notification struct {
	ID string `json:"id,omitempty"`
	TimestampMillis   uint64 `json:"timestampMillis"`
	Destination Destination `json:"destination"`
	Acknowledged bool `json:"acknowledged"`
	AlertUID string `json:"alertUID,omitempty"`
	LogUID string `json:"logUID,omitempty"`
	Body map[string]string `json:"body"`

	// The notification profile
	Profile string `json:"profile"`
	Labels map[string]string `json:"labels"`
	DocumentationURL string `json:"documentationURL"`
	SnoozeURL string `json:"snoozeURL"`
}
