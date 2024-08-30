package grouping

import (
	"testing"
)

func TestInitRules(t *testing.T) {
	rules := []*Grouping{
		{
			If:      `identity.kind == "host"`,
			GroupBy: []string{"identity.hostname", "identity.process"},
		},
	}
	Startup(rules)
}
