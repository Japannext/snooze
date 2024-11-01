package models

const SNOOZE_INDEX = "v2-snoozes"

type Snooze struct {
	Base
	GroupName string `json:"groupName"`
	Hash string `json:"hash"`
	Reason string `json:"reason"`
	StartAt Time `json:"startAt"`
	ExpireAt Time `json:"expireAt"`
}
