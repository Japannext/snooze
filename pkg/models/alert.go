package models

const (
	ActiveAlertIndex = "v2-active-alerts"
	AlertHistoryIndex = "v2-alerts-history"
)


type AlertBase struct {
	Base

	Hash string `json:"hash"`

	Source Source `json:"source"`

	Identity map[string]string `json:"identity"`

	// Timestamps
	StartsAt Time `json:"startsAt"`

	// Name of the alert
	AlertName string `json:"alertName"`
	// Group of the alert
	AlertGroup string `json:"alertGroup"`

	// Text representing the severity
	SeverityText string `json:"severityText,omitempty"`
	// Number representing the severity. Useful for filters (severity higher than a given value)
	SeverityNumber int32 `json:"severityNumber,omitempty"`

	TraceID string `json:"traceID,omitempty"`

	// Text representing the severity
	// SeverityText string `json:"severityText,omitempty"`
	// Number representing the severity. Useful for filters (severity higher than a given value)
	// SeverityNumber int32 `json:"severityNumber,omitempty"`

	Description string `json:"description"`
	Summary     string `json:"summary"`

	// Key-value representing the main resource identifiers.
	// Examples: host, pod, disk name, IP address
	Labels map[string]string `json:"labels,omitempty"`
}

type ActiveAlert struct {
	AlertBase

	LastHit Time `json:"lastHit"`
}

type AlertRecord struct {
	AlertBase

	EndsAt Time `json:"endsAt"`
}

type AlertStatus struct {
	ID        string `json:"id"`
	LastCheck Time   `json:"lastCheck"`
	NextCheck Time   `json:"nextCheck"`
}

func (item *ActiveAlert) Context() map[string]interface{} {
	return map[string]interface{}{
		"source":      map[string]string{"kind": item.Source.Kind, "name": item.Source.Name},
		"identity":    item.Identity,
		"labels":      item.Labels,
		"summary":     item.Summary,
		"description": item.Description,
		"alertName":   item.AlertName,
		"alertGroup":  item.AlertGroup,
	}
}

type AlertUpdate struct {
	Document *ActiveAlert `json:"doc"`
	Upsert   *ActiveAlert `json:"upsert"`
}

func init() {
	OpensearchIndices[ActiveAlertIndex] = Indice{
		Settings: IndexSettings{1, 2},
		Mappings: IndexMapping{
			Properties: map[string]MappingProps{
				"hash":        {Type: "keyword"},
				"startsAt":    {Type: "date", Format: "epoch_millis"},
				"source.kind": {Type: "keyword"},
				"source.name": {Type: "keyword"},
				"labels":      {Type: "object"},
			},
		},
	}
	OpensearchIndexTemplates[AlertHistoryIndex] = IndexTemplate{
		Version: 1,
		IndexPatterns: []string{AlertHistoryIndex},
		DataStream:    map[string]map[string]string{"timestamp_field": {"name": "startsAt"}},
		Template: Indice{
			Settings: IndexSettings{1, 2},
			Mappings: IndexMapping{
				Properties: map[string]MappingProps{
					"hash":        {Type: "keyword"},
					"startsAt":    {Type: "date", Format: "epoch_millis"},
					"endsAt":      {Type: "date", Format: "epoch_millis"},
					"source.kind": {Type: "keyword"},
					"source.name": {Type: "keyword"},
					"labels":      {Type: "object"},
				},
			},
		},
	}
}
