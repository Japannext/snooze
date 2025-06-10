package models

// A selection is something that satisfied a Selector.
type Selection struct {
	Hash          string            `json:"hash"`
	MatchLabels   map[string]string `json:"matchLabels,omitempty"`
	MatchIdentity map[string]string `json:"matchIdentity,omitempty"`
}
