package grouping

import (
	"testing"
)

func TestInitRules(t *testing.T) {
	rules := []*Rule{
		{
			If:      `alert.source.kind == "syslog"`,
			GroupBy: []string{"$.labels.host", "$.labels.process"},
		},
	}
	InitRules(rules)
}
