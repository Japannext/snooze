package models

const GROUP_INDEX = "v2-groups"

type Group struct {
	Name string `json:"name"`
	// Hash value of the key-value
	Hash string `json:"hash"`
	// Human readable information about the group
	Labels map[string]string `json:"labels,omitempty"`
}
