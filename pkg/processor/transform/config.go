package transform

import (
	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/lang"
)

type Transform struct {
	Name	 string		`yaml:"name"`
	If       string     `yaml:"if"`

	// Transforms []transformConfig `yaml:"transforms"`
	// Transform transformConfig `yaml:",inline"`

	internal struct {
		condition *lang.Condition
		transform Transformation
	}
}

type transformConfig struct {
	Set *SetTransform `yaml:"set"`
}

func (tr *Transform) Load() {
	if tr.If != "" {
		condition, err := lang.NewCondition(tr.If)
		if err != nil {
			log.Fatal(err)
		}
		tr.internal.condition = condition
	}

	/*
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
	*/
	// TODO
}

type Transformation interface {
	Load()
	Process(*models.Log) error
}
