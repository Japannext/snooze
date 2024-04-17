package v2

type Destination struct {
  // Type of destination (e.g mail/googlechat/patlite/browser)
  Kind string `json:"kind"`
  // Name of the instance of destination backend (e.g. prod/example.com)
  Name string `json:"name"`
}

type Notification struct {
  Destination Destination `json:"destination"`

  Timestamp uint64              `json:"timestamp"`
  ObservedTimestamp uint64      `json:"observed_timestamp,omitempty"`
}
