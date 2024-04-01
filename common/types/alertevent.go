package types

type SeverityNumber = int32
type KeyValue = map[string]string

const (
  LogKind = "log"
  NotificationKind = "notification"
  RatelimitKind = "ratelimit"
)

type AlertEventV2 struct {
  Kind string                   `json:"kind" cql:"kind"`
  Timestamp uint64              `json:"timestamp" cql:"timestamp"`
  GroupHash []byte              `json:"group_hash,omitempty" cql:"group_hash"`
  GroupKv KeyValue              `json:"group_kv,omitempty" cql:"group_kv"`
  SeverityText string           `json:"severity_text,omitempty" cql:"severity_text"`
  SeverityNumber SeverityNumber `json:"severity_number,omitempty" cql:"severity_nb"`
  Resource KeyValue             `json:"resource,omitempty" cql:"resource"`
  Attributes KeyValue           `json:"attributes,omitempty" cql:"attributes"`
  Body KeyValue                 `json:"body" cql:"body"`
}
