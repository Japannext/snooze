package transform

import (
	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/lang"
)

type Transform struct {
	Name	 string		`yaml:"name"`
	If       string     `yaml:"if"`

	Parent *ParentTransform `yaml:",inline"`
	Set *SetTransform `yaml:"set,omitempty"`

	internal struct {
		condition *lang.Condition
		transform Transformation
	}
}

func (tr *Transform) Load() {
	if tr.If != "" {
		condition, err := lang.NewCondition(tr.If)
		if err != nil {
			log.Fatal(err)
		}
		tr.internal.condition = condition
	}

	switch {
		case tr.Parent != nil:
			tr.Parent.Load()
			tr.internal.transform = tr.Parent
		case tr.Set != nil:
			tr.Set.Load()
			tr.internal.transform = tr.Set
		default:
			log.Fatalf("transform `%s` is empty", tr.Name)
	}
}

type Transformation interface {
	Load()
	Process(*models.Log) error
}
