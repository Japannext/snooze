package condition

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/japannext/snooze/common/field"
)

func f(name, key string) field.AlertField {
	return field.AlertField{Name: name, SubKey: key}
}

func eq(name, key, v string) *Condition {
	return &Condition{Kind: "equal", Equal: &EqualCondition{f(name, key), v}}
}
func noteq(name, key, v string) *Condition {
	return &Condition{Kind: "not_equal", NotEqual: &NotEqualCondition{f(name, key), v}}
}
func tmatch(name, key, v string) *Condition {
	return &Condition{Kind: "match", Match: &MatchCondition{f(name, key), v, nil}}
}

func TestParser(t *testing.T) {
	var tests = []struct {
		Text     string
		Expected *Condition
	}{
		{`labels[host.name] == "host01"`, eq("labels", "host.name", "host01")},
		{`labels[host.name] != "host01"`, noteq("labels", "host.name", "host01")},
		{`labels[host.name] =~ "host"`, tmatch("labels", "host.name", "host")},
	}

	for _, data := range tests {
		t.Run(data.Text, func(t *testing.T) {
			c, err := Parse(data.Text)
			assert.NoError(t, err)
			assert.Equal(t, data.Expected, c)
		})
	}

}
