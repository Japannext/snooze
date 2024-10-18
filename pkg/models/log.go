package models

const LOG_INDEX = "v2-logs"

type Log struct {

	ID string `json:"id,omitempty"`

	// Several timestamps important to the log.
	Timestamp Timestamp `json:"timestamp"`

	// Information regarding the source plugin that created the log.
	Source Source `json:"source"`

	// What server/container/pod emitted the log.
	Identity map[string]string `json:"identity,omitempty"`

	Profile string `json:"profile,omitempty"`
	Pattern string `json:"pattern,omitempty"`

	// ID for opentelemetry traces
	TraceID string `json:"traceID"`

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

	// Indicate that this message is part of an active check, and that
	// this is the URL of the callback
	ActiveCheckURL string `json:"activeCheckURL,omitempty"`

	// In case the log encountered a live error during processing, this
	// can be indicated here.
	Error string `json:"error,omitempty"`

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
		"timestamp": item.Timestamp,
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

func init() {
	index := IndexTemplate{
		Version: 0,
		IndexPatterns: []string{LOG_INDEX},
		DataStream: map[string]map[string]string{"timestamp_field": {"name": "timestamp.display"}},
		Template: Indice{
			Settings: IndexSettings{1, 2},
			Mappings: IndexMapping{
				Properties: map[string]MappingProps{
					"timestamp.display": {Type: "date", Format: "epoch_millis"},
					"timestamp.actual": {Type: "date", Format: "epoch_millis"},
					"timestamp.observed": {Type: "date", Format: "epoch_millis"},
					"timestamp.processed": {Type: "date", Format: "epoch_millis"},
					"timestamp.warning": {Type: "keyword"},
					"source.kind": {Type: "keyword"},
					"source.name": {Type: "keyword"},
					"identity": {Type: "object"},
					//"group.hash":   {Type: "keyword"},
					//"group.labels": {Type: "object"},
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
