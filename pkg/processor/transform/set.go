package transform

import (
	"github.com/japannext/snooze/pkg/models"
)

type SetTransform struct {
	Target string `yaml:"target" json:"target"`
	Value  string `yaml:"value" json:"value"`
}

func (tr *SetTransform) Load() {}

func (tr *SetTransform) Process(item *models.Log) error {
	return nil
}
