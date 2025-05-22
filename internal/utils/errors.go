package utils

import (
	"errors"
	"strings"
)

type MultiError struct {
	header string
	errors []error
}

func NewMultiError(header string) *MultiError {
	return &MultiError{
		header: header,
		errors: []error{},
	}
}

func (merr *MultiError) Append(s string) {
	merr.errors = append(merr.errors, errors.New(s))
}

func (merr *MultiError) AppendErr(err error) {
	merr.errors = append(merr.errors, err)
}

func (merr *MultiError) Error() string {
	var s strings.Builder
	s.WriteString(merr.header)
	s.WriteString("\n")
	for _, err := range merr.errors {
		s.WriteString(err.Error())
		s.WriteString("\n")
	}
	return s.String()
}

func (merr *MultiError) HasErrors() bool {
	return (len(merr.errors) != 0)
}
