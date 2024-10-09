package models

type Source struct {
	// Source kind/protocol (e.g. syslog, OTEL, prometheus, etc)
	Kind string `json:"kind"`
	// The source instance name (e.g. "prod-relay", "host01", "tenant-A")
	Name string `json:"name,omitempty"`
}
