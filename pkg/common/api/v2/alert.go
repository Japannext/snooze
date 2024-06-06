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

type Source struct {
	// Source kind/protocol (e.g. syslog, OTEL, prometheus, etc)
	Kind string `json:"kind"`
	// The source instance name (e.g. "prod-relay", "host01", "tenant-A")
	Name string `json:"name"`
}

type Mute struct {
	// Enable the muting
	Enabled bool `json:"enabled"`
	// The reason it was muted. `snooze`/`silence`/`test`
	Component string `json:"component"`
	// Name of the silence rule / snooze rule that muted the alert
	Rule string `json:"rule"`
	// Skip the notification step. Usually on.
	SkipNotification bool `json:"skipNotification"`
	// Skip storing into the database (opensearch). Usually used for testing
	SkipStorage bool `json:"skipStorage"`
	// A test alert, which will not trigger anything. Mainly used for internal metrics
	// and active monitoring of the snooze pipelines.
	SilentTest bool

	// A test manually performed by a human. Will trigger everything
	// normally (patlite, etc), but will be marked as such in the web interface.
	HumanTest bool
}

// A function that can match an alert
type AlertMatcher = func(Alert) bool

// A function that can modify an alert in-place
type AlertModifier = func(Alert)
