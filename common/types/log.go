package types

import (
  "github.com/scylladb/gocqlx/v2/table"
)

type SeverityNumber = int32
type KeyValue = map[string]string

const (
  LogKind = "log"
  NotificationKind = "notification"
  RatelimitKind = "ratelimit"
)

type LogV2 struct {
  Kind string `json:"kind"`
  Timestamp uint64 `json:"timestamp"`
  GroupHash []byte `json:"groupHash,omitempty"`
  GroupKv KeyValue `json:"groupKv,omitempty"`
  SeverityText string `json:"severityText,omitempty"`
  SeverityNumber SeverityNumber `json:"severityNumber,omitempty"`
  Resource KeyValue `json:"resource,omitempty"`
  Attributes KeyValue `json:"attributes,omitempty"`
  Body KeyValue `json:"body"`
}

var LogV2Metadata = table.Metadata{
    Name: "log",
    Columns: []string{"timestamp", "kind", "groupHash", "groupKv", "resource", "attributes", "body"},
    PartKey: []string{"groupHash"},
    SortKey: []string{"timestamp"},
}
