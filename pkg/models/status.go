package models

type LogStatus struct {
	Kind     int    `json:"kind"`
	ObjectID string `json:"objectID,omitempty"`
	Reason   string `json:"reason"`

	SkipNotification bool `json:"skipNotification"`
	SkipStorage      bool `json:"skipStorage"`
}

type LogStatusKind int

const (
	LogActive int = iota
	LogSnoozed
	LogSilenced
	LogRatelimited
	LogDropped
	LogActiveCheck
	LogAcked
)

func (status *LogStatus) Change(kind int) bool {
	// Either more important, or if same importance,
	// the first wins
	if kind > status.Kind {
		status.Kind = kind
		status.ObjectID = ""
		status.Reason = ""
		return true
	}
	return false
}

type AlertStatus2 struct {
	Kind int `json:"kind"`
}

const (
	AlertFiring = iota
	AlertResolved
	AlertUnknown
	AlertAcked
)
