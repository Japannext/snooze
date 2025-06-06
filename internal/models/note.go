package models

// A note is a message written by a human or machine,
// that is attached to a problem object. It is displayed
// in a timeline alongside logs, and alert state change.
type Note struct {
	Base
	ProblemID string     `json:"problemID"`
	Kind      string     `json:"kind"`
	Username  string     `json:"username"`
	Content   string     `json:"content"`
	Reactions []Reaction `json:"reactions"`
}

// A reaction to an event (like 👍)
type Reaction struct {
	Username string `json:"username"`
	Emoji    string `json:"emoji"`
}
