package v2

type Destination struct {
	// Type of destination (e.g mail/googlechat/patlite/browser)
	Kind string `json:"kind"`
	// Name of the instance of destination backend (e.g. prod/example.com)
	Name string `json:"name"`
}

type Notification struct {
	TimestampMillis   uint64 `json:"timestampMillis"`
	Destination Destination `json:"destination"`
	AlertUID string `json:"alertUID,omitempty"`
	LogUID string `json:"logUID,omitempty"`
	Body map[string]string `json:"body"`
}
