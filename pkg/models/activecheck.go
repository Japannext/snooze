package models

type ActiveCheck struct {
	// The destination it ended routed as
	Destination Destination `json:"destination"`
	// In case the log encountered a live error in processing or source, the
	// error can be non-null and indicated here.
	Error string `json:"error,omitempty"`
}

