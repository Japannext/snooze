package models

const ACK_INDEX = "v2-ack"

// An acknowledgement given to a log or
// group of logs.
type Ack struct {
	Base
	Time Time `json:"timestamp"`
	Username string `json:"username"`
	Reason string `json:"reason"`
	LogIDs []string `json:"logIDs"`
}
