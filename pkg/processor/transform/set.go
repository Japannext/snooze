package transform

import (
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

type SetRule struct {
	Target string `yaml:"target" json:"target"`
	Value  string `yaml:"value" json:"value"`
}

func (rule *SetRule) Compute() Interface {
	return nil
}

type computedSet struct {
}

func (s *computedSet) Process(alert *api.Alert) error {
	return nil
}
