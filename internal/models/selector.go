package models

import (
	"fmt"
)

// Selector for labels. A hash may be calculated and added to the object.
type Selector struct {
	Hash          string            `json:"hash,omitempty"`
	MatchLabels   map[string]string `json:"matchLabels,omitempty" yaml:"matchLabels"`
	MatchIdentity map[string]string `json:"matchIdentity,omitempty" yaml:"matchIdentity"`
}

// Validate the object.
func (s *Selector) Validate() error {
	if len(s.MatchLabels) > 0 {
		return nil
	}
	if len(s.MatchIdentity) > 0 {
		return nil
	}
	return fmt.Errorf("empty selector")
}
