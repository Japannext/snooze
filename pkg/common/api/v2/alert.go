package v2

import (
	"strings"
)

type Alert struct {
	Source Source `json:"source"`

	Timestamp         uint64 `json:"timestamp"`
	ObservedTimestamp uint64 `json:"observedTimestamp,omitempty"`

	GroupHash   []byte            `json:"groupHash,omitempty"`
	GroupLabels map[string]string `json:"groupLabels,omitempty"`

	// Text representing the severity
	SeverityText string `json:"severityText,omitempty"`
	// Number representing the severity. Useful for filters (severity higher than a given value)
	SeverityNumber int32 `json:"severityNumber,omitempty"`

	// Key-value representing the main resource identifiers.
	// Examples: host, pod, disk name, IP address
	Labels map[string]string `json:"labels,omitempty"`
	// Additional attributes that may not be essential to identify
	// the resource the alert is concerned.
	// Examples: CPU arch, OS version, program version, URLs to source/documentation
	Attributes map[string]string `json:"attributes,omitempty"`
	// The main body of the alert. This represents what should be read
	// by operators to act on it.
	// Examples: log message, prometheus alert summary/details/value
	Body map[string]string `json:"body"`

	// Mute the alert. This may skip notifications, or skip even display
	Mute Mute `json:"mute"`
}

func (a *Alert) String() string {
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
	for k, v := range a.Body {
		s.WriteString(k)
		s.WriteString(":")
		s.WriteString(v)
		s.WriteString("; ")
	}

	return s.String()
}
