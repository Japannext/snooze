package v2

type AlertState struct {
  Kind string                   `json:"kind"`
  Name string                   `json:"name"`
  State string                  `json:"state"`

  Timestamp uint64              `json:"timestamp"`
  GroupHash []byte              `json:"group_hash,omitempty"`
  GroupKv KeyValue              `json:"group_kv,omitempty"`
  NotificationChannels []string `json:"notification_channels"`
  SeverityText string           `json:"severity_text,omitempty"`
  SeverityNumber SeverityNumber `json:"severity_number,omitempty"`
  Resource KeyValue             `json:"resource,omitempty"`
  Attributes KeyValue           `json:"attributes,omitempty"`
  Body KeyValue                 `json:"body"`
}
