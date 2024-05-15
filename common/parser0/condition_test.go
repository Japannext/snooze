package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/japannext/snooze/common/condition"
	"github.com/japannext/snooze/common/field"
)

func TestParser(t *testing.T) {
	_f := func(name, subkey string) *field.AlertField {
		return &field.AlertField{Name: name, SubKey: subkey}
	}
	var tests = []struct {
		Text     string
		Expected *condition.Condition
	}{
		{`labels[host.name] = "host01"`, condition.NewEqual(_f("labels", "host.name"), "host01")},
		{`labels[host.name] != "host01"`, condition.NewNotEqual(_f("labels", "host.name"), "host01")},
		{`labels[host.name] =~ "host"`, condition.NewMatch(_f("labels", "host.name"), "host")},
		{`has labels[service.name]`, condition.NewHas(_f("labels", "service.name"))},
		{`(has labels[service.name]) and (has labels[environment])`, condition.NewAnd(
			condition.NewHas(_f("labels", "service.name")),
			condition.NewHas(_f("labels", "environment")),
		)},
	}

	for _, data := range tests {
		t.Run(data.Text, func(t *testing.T) {
			fmt.Printf("Running %s\n", data.Text)
			c, err := ParseCondition(data.Text)
			assert.NoError(t, err)
			assert.Equal(t, data.Expected, c)
		})
	}

}
