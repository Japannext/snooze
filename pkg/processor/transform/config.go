package transform

import (
	"context"

	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/models"
)

type Transform struct {
	Name string `json:"name"         yaml:"name"`
	If   string `json:"if,omitempty" yaml:"if"`

	Actions []Action `json:"actions" yaml:"actions"`

	internal struct {
		condition *lang.Condition
		actions   []Transformation
	}
}

type Transformation interface {
	Process(context.Context, *models.Log) (context.Context, error)
}

type Action struct {
	Set   *SetAction   `json:"set,omitempty"   yaml:"set"`
	Unset *UnsetAction `json:"unset,omitempty" yaml:"unset"`
	Regex *RegexAction `json:"regex,omitempty" yaml:"regex"`
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
