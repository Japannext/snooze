package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/japannext/snooze/common/field"
)

func TestParseField(t *testing.T) {
	_f := func(name, subkey string) *field.AlertField {
		return &field.AlertField{Name: name, SubKey: subkey}
	}

	tests := []struct {
		text     string
		expected *field.AlertField
		success  bool
	}{
		{"severity", _f("severity", ""), true},
		{"labels[service.name]", _f("labels", "service.name"), true},
		{"attributes[appname]", _f("attributes", "appname"), true},
		{"labels[test", nil, false},
	}
	for _, e := range tests {
		t.Run(e.text, func(t *testing.T) {
			res, err := ParseField(e.text)
			if e.success {
				assert.NoError(t, err)
				assert.Equal(t, e.expected, res)
			} else {
				assert.Error(t, err)
				fmt.Printf("[ERROR] %s", err)
			}
		})
	}
}
