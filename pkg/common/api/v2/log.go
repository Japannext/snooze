package v2

import (
	"strings"
)

type ProcessedLog struct {
	Log *Log `json:",inline"`

	LogPattern string `json:"logPattern"`
	Identity map[string]string `json:"identity"`
	Group map[string]string `json:"group"`
}

type Log struct {
	ID string `json:"_id,omitempty"`

	// The type of alert (syslog, opentelemetry, snmptrap, prometheus rule...),
	// as well as the name of the instance (if any)
	Source Source `json:"source"`

	Identity map[string]string `json:"identity,omitempty"`

	TimestampMillis         uint64 `json:"timestampMillis"`
	ObservedTimestampMillis uint64 `json:"observedTimestampMillis,omitempty"`

	Profile string `json:"profile,omitempty"`
	Pattern string `json:"pattern,omitempty"`

	GroupHash   []byte            `json:"groupHash,omitempty"`
	GroupLabels map[string]string `json:"groupLabels,omitempty"`

	// Text representing the severity
	SeverityText string `json:"severityText,omitempty"`
	// Number representing the severity. Useful for filters (severity higher than a given value)
	SeverityNumber int32 `json:"severityNumber,omitempty"`

	// Key-value representing the main resource identifiers.
	// Examples: host, pod, disk name, IP address
	Labels map[string]string `json:"labels,omitempty"`

	// The main message of the log
	Message string `json:"message"`

	// Mute the alert. This may skip notifications, or skip even display
	Mute Mute `json:"mute,omitempty"`
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
