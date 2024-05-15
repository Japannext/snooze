package transform

import (
	api "github.com/japannext/snooze/common/api/v2"
	"github.com/japannext/snooze/common/field"
	"github.com/japannext/snooze/common/parser"
)

type SetRule struct {
	Target string `yaml:"target" json:"target"`
	Value  string `yaml:"value" json:"value"`
}

func (rule *SetRule) Compute() Interface {
	f, err := parser.ParseField(rule.Target)
	if err != nil {
		log.Fatal(err)
	}
	return &computedSet{*f, rule.Value}
}

type computedSet struct {
	TargetField field.AlertField
	Value       string
}

func (s *computedSet) Process(alert *api.Alert) error {
	s.TargetField.Set(alert, s.Value)
	return nil
}
