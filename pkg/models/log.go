package models

const LOG_INDEX = "v2-logs"

type Log struct {
	Base

	// The timestamp when the event actually happened (when known)
	ActualTime Time `json:"actualTime"`
	// The timestamp when the event entered the system (snooze or otel)
	ObservedTime Time `json:"observedTime,omitempty"`
	// TIme used in web interfaces for sorting. Equals to the value
	// that makes most sense between observedTime and actualTime.
	DisplayTime Time `json:"displayTime"`

	// Information regarding the source plugin that created the log.
	Source Source `json:"source"`

	// What server/container/pod emitted the log.
	Identity map[string]string `json:"identity,omitempty"`

	Profile string `json:"profile,omitempty"`
	Pattern string `json:"pattern,omitempty"`

	// ID for opentelemetry traces
	TraceID string `json:"traceID"`

	// Groups used for search and process
	Groups []*Group `json:"groups"`

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

	// Indicate that this message is part of an active check, and that
	// this is the URL of the callback
	ActiveCheckURL string `json:"activeCheckURL,omitempty"`

	// In case the log encountered a live error during processing, this
	// can be indicated here.
	Error string `json:"error,omitempty"`

	Status Status `json:"status"`
}

// Used by template systems in transforms/profiles/etc
func (item *Log) Context() map[string]interface{} {
	return map[string]interface{}{
		"actualTime": item.ActualTime,
		"source": item.Source,
		"identity": item.Identity,
		"labels": item.Labels,
		"message": item.Message,
	}
}

type Process struct {
	Profile string `json:"profile,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}

func init() {
	index := IndexTemplate{
		Version: 2,
		IndexPatterns: []string{LOG_INDEX},
		DataStream: map[string]map[string]string{"timestamp_field": {"name": "displayTime"}},
		Template: Indice{
			Settings: IndexSettings{1, 2},
			Mappings: IndexMapping{
				Properties: map[string]MappingProps{
					"displayTime": {Type: "date", Format: "epoch_millis"},
					"actualTime": {Type: "date", Format: "epoch_millis"},
					"observedTime": {Type: "date", Format: "epoch_millis"},
					"groups.name": {Type: "keyword"},
					"groups.labels": {Type: "object"},
					"groups.hash": {Type: "keyword"},
					"source.kind": {Type: "keyword"},
					"source.name": {Type: "keyword"},
					"identity": {Type: "object"},
					"profile": {Type: "keyword"},
					"pattern": {Type: "keyword"},
					"labels":      {Type: "object"},
					"message":        {Type: "text"},
					"mute.skipNotification": {Type: "boolean"},
					"mute.skipStorage": {Type: "boolean"},
					"mute.reason": {Type: "keyword"},
				},
			},
		},
	}
	INDEXES = append(INDEXES, index)
}
