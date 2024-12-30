package transform

import (
	"context"

	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/models"
)

type Transform struct {
	Name string `yaml:"name" json:"name"`
	If   string `yaml:"if" json:"if,omitempty"`

	Actions []Action `yaml:"actions" json:"actions"`

	internal struct {
		condition *lang.Condition
		actions   []Transformation
	}
}

type Transformation interface {
	Process(context.Context, *models.Log) (context.Context, error)
}

type Action struct {
	Set   *SetAction   `yaml:"set" json:"set,omitempty"`
	Unset *UnsetAction `yaml:"unset" json:"unset,omitempty"`
	Regex *RegexAction `yaml:"regex" json:"regex,omitempty"`
}

func LoadActions(actions []Action) []Transformation {
	var transforms []Transformation
	for _, action := range actions {
		switch {
		case action.Set != nil:
			transforms = append(transforms, action.Set.Load())
		case action.Unset != nil:
			transforms = append(transforms, action.Unset.Load())
		case action.Regex != nil:
			transforms = append(transforms, action.Regex.Load())
		default:
			log.Fatalf("action is empty")
		}
	}
	return transforms
}

func (tr *Transform) Load() {
	if tr.If != "" {
		condition, err := lang.NewCondition(tr.If)
		if err != nil {
			log.Fatal(err)
		}
		tr.internal.condition = condition
	}

	tr.internal.actions = LoadActions(tr.Actions)
}
