package v2

type Alert struct {
  Source struct {
    // Source kind/protocol (e.g. syslog, OTEL, prometheus, etc)
    Kind string `json:"kind"`
    // The source instance name (e.g. "prod-relay", "host01", "tenant-A")
    Name string `json:"name"`
  } `json:"source"`

  Timestamp uint64              `json:"timestamp"`
  ObservedTimestamp uint64      `json:"observedTimestamp,omitempty"`

  GroupHash []byte              `json:"groupHash,omitempty"`
  GroupLabels map[string]string     `json:"groupLabels,omitempty"`

  // Text representing the severity
  SeverityText string           `json:"severityText,omitempty"`
  // Number representing the severity. Useful for filters (severity higher than a given value)
  SeverityNumber int32          `json:"severityNumber,omitempty"`

  // Key-value representing the main resource identifiers.
  // Examples: host, pod, disk name, IP address
  Labels map[string]string      `json:"labels,omitempty"`
  // Additional attributes that may not be essential to identify
  // the resource the alert is concerned.
  // Examples: CPU arch, OS version, program version, URLs to source/documentation
  Attributes map[string]string  `json:"attributes,omitempty"`
  // The main body of the alert. This represents what should be read
  // by operators to act on it.
  // Examples: log message, prometheus alert summary/details/value
  Body map[string]string        `json:"body"`

  // A test alert, which will not trigger anything. Mainly used for internal metrics
  // and active monitoring of the snooze pipelines.
  SilentTest bool

  // A test manually performed by a human. Will trigger everything
  // normally (patlite, etc), but will be marked as such in the web interface.
  HumanTest bool

  // Set when the alert was snoozed
  Snoozed bool `json:"snoozed"`

  // Set when the alert was snooze. Correspond to the name of the rule
  // that snoozed the alert.
  SnoozeRuleName string `json:"snoozeRuleName"`
}
