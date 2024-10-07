package models

type ActiveCheckCallback struct {
	// Delay in milliseconds between the active check being emitted and it
	// being processed by the system.
	DelayMillis uint64 `json:"delayMillis"`
	// In case the log encountered a live error in processing or source, the
	// error can be non-null and indicated here.
	Error string `json:"error,omitempty"`
}

