package utils

import (
	"fmt"
)

// A checker that check duplicate names.
type DuplicateChecker struct {
	mapper map[string]bool
}

func NewDuplicateChecker() *DuplicateChecker {
	return &DuplicateChecker{mapper: map[string]bool{}}
}

func (checker *DuplicateChecker) Check(name string) error {
	if _, ok := checker.mapper[name]; ok {
		return fmt.Errorf("duplicate entry named `%s`", name)
	}

	checker.mapper[name] = true

	return nil
}
