package models

type Status struct {
	Kind string `json:"kind"`
	ObjectID string `json:"objectID,omitempty"`
	Reason string `json:"reason"`

	SkipNotification bool `json:"skipNotification"`
	SkipStorage bool `json:"skipStorage"`
}

type LogStatus struct {
	Kind LogStatusKind `json:"kind"`
	ObjectID string `json:"objectID,omitempty"`
	Reason string `json:"reason"`

	SkipNotification bool `json:"skipNotification"`
	SkipStorage bool `json:"skipStorage"`
}

type LogStatusKind int

const (
	LogActive LogStatusKind = iota
	LogSnoozed
	LogSilenced
	LogRatelimited
	LogDropped
	LogActiveCheck
	LogAcked
)

func (kind LogStatusKind) String() string {
	return [...]string{
		"Active",
		"Snoozed",
		"Silenced",
		"Ratelimited",
		"Dropped",
		"ActiveCheck",
		"Acked",
	}[kind]
}

func (status *LogStatus) Change(kind LogStatusKind) bool {
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
