package transform

import (
	"github.com/japannext/snooze/pkg/models"
)

type ParentTransform struct {
	Children []*Transform `yaml:"children"`
}

func (parent *ParentTransform) Load() {
	for _, tr := range parent.Children {
		tr.Load()
	}
}

func (parent *ParentTransform) Process(item *models.Log) error {
	for _, tr := range parent.Children {
		err := tr.internal.transform.Process(item)
		if err != nil {
			return err
		}
	}
	return nil
}
