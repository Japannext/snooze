package utils

import (
	"fmt"
)

type NameValidator struct {
	mapper   map[string]bool
	notEmpty bool
}

func NewNameValidator(notEmpty bool) *NameValidator {
	return &NameValidator{mapper: map[string]bool{}, notEmpty: notEmpty}
}

func (validator *NameValidator) Check(name string) error {
	if validator.notEmpty && name == "" {
		return fmt.Errorf("found entry with empty `name`")
	}
	if _, ok := validator.mapper[name]; ok {
		return fmt.Errorf("duplicate entry named `%s`", name)
	}
	validator.mapper[name] = true
	return nil
}
