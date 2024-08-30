package transform

import (
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

type SetTransform struct {
	Target string `yaml:"target" json:"target"`
	Value  string `yaml:"value" json:"value"`
}

func (tr *SetTransform) Load() {}

func (tr *SetTransform) Process(item *api.Log) error {
	return nil
}
