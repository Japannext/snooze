package models

type SourceActiveCheck struct {
	// In case the log encountered a live error in processing or source, the
	// error can be non-null and indicated here.
	Error string `json:"error,omitempty"`
}

type DestinationActiveCheck struct {
	Destination Destination `json:"destination"`
	Error       string      `json:"error,omitempty"`
}
