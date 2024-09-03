package v2

import (
	"strings"
)

type Log struct {
	ID string `json:"id,omitempty"`

	TimestampMillis         uint64 `json:"timestampMillis"`
	ObservedTimestampMillis uint64 `json:"observedTimestampMillis,omitempty"`

	// Information regarding the source plugin that created the log.
	Source Source `json:"source"`

	// What server/container/pod emitted the log.
	Identity map[string]string `json:"identity,omitempty"`

	Profile string `json:"profile,omitempty"`
	Pattern string `json:"pattern,omitempty"`

	// Grouping
	Group LogGroup `json:"group"`

	// Text representing the severity
	SeverityText string `json:"severityText,omitempty"`
	// Number representing the severity. Useful for filters (severity higher than a given value)
	SeverityNumber int32 `json:"severityNumber,omitempty"`

	// Key-value representing the main resource identifiers.
	// Examples: host, pod, disk name, IP address
	Labels map[string]string `json:"labels,omitempty"`

	// The main message of the log
	Message string `json:"message"`

	// Details written during snooze-process
	Process *Process `json:"process,omitempty"`

	// Mute the alert. This may skip notifications, or skip even display
	Mute Mute `json:"mute,omitempty"`
}

func (item *Log) GetID() string { return item.ID }
func (item *Log) SetID(id string) { item.ID = id }

type HasContext interface {
	Context() map[string]interface{}
}

// Used by template systems in transforms/profiles/etc
func (item *Log) Context() map[string]interface{} {
	return map[string]interface{}{
		"timestamp": item.TimestampMillis,
		"source": item.Source,
		"identity": item.Identity,
		"labels": item.Labels,
		"message": item.Message,
	}
}

type Process struct {
	Profile string `json:"profile,omitempty"`
	Pattern string `json:"pattern,omitempty"`
	Group LogGroup `json:"group"`
}
// A group of logs, that can be uniquely identified by a hash
type LogGroup struct {
	Hash string `json:"hash,omitempty"`
	Labels map[string]string `json:"labels"`
}

type LogResults struct {
	Items []Log `json:"items"`
	// Indicate the total number of logs during the requested time period
	Total int `json:"total"`
	// Indicate if the total is exact, or if it reached
	// the backend limit per response.
	More bool `json:"more"`
}

func (a *Log) String() string {
	var s strings.Builder

	// Source
	s.WriteString("[")
	s.WriteString(a.Source.Kind)
	if a.Source.Name != "" {
		s.WriteString("/")
		s.WriteString(a.Source.Name)
	}
	s.WriteString("] ")

	// Labels
	s.WriteString("[")
	for k, v := range a.Labels {
		s.WriteString(k)
		s.WriteString("=")
		s.WriteString(v)
		s.WriteString(", ")
	}
	s.WriteString("] ")

	// Body
	s.WriteString(a.Message)

	return s.String()
}

type Source struct {
	// Source kind/protocol (e.g. syslog, OTEL, prometheus, etc)
	Kind string `json:"kind"`
	// The source instance name (e.g. "prod-relay", "host01", "tenant-A")
	Name string `json:"name,omitempty"`
}
