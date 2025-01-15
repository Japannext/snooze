package transform

import (
	"fmt"

	"github.com/japannext/snooze/pkg/common/lang"
)

type Processor struct {
	transforms []*transform
}

type Config struct {
	Transforms []*TrConfig `json:"transforms" yaml:"transforms"`
}

type TrConfig struct {
	Name     string         `json:"name"         yaml:"name"`
	If       string         `json:"if,omitempty" yaml:"if"`
	Actions  []*ActionConfig `json:"actions"     yaml:"actions"`
}

func New(cfg Config) (*Processor, error) {
	p := &Processor{}

	for _, trConfig := range cfg.Transforms {
		tr, err := newTransform(trConfig)
		if err != nil {
			return p, fmt.Errorf("in transform '%s': %w", trConfig.Name, err)
		}

		p.transforms = append(p.transforms, tr)
	}

	return p, nil
}

type transform struct {
	cfg *TrConfig

	condition *lang.Condition
	actions []ActionInterface
}

func newTransform(cfg *TrConfig) (*transform, error) {
	tr := &transform{cfg: cfg}

	if cfg.If != "" {
		condition, err := lang.NewCondition(cfg.If)
		if err != nil {
			return tr, fmt.Errorf("error in condition `%s`: %w", cfg.If, err)
		}
		tr.condition = condition
	}

	for index, actionConfig := range cfg.Actions {
		action, err := NewAction(actionConfig)
		if err != nil {
			return tr, fmt.Errorf("in action #%d: %w", index+1, err)
		}
		tr.actions = append(tr.actions, action)
	}

	return tr, nil
}
