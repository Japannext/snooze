package models

//
type Status struct {
	Kind string `json:"kind"`
	ObjectID string `json:"objectID,omitempty"`
	Reason string `json:"reason"`

	SkipNotification bool `json:"skipNotification"`
	SkipStorage bool `json:"skipStorage"`
}
