package transform

import (
	"fmt"

	"github.com/japannext/snooze/internal/lang"
	"github.com/japannext/snooze/internal/pipeline"
)

type Processor struct {
	transforms []*transform
}

func New(trs []pipeline.Transform) (*Processor, error) {
	p := &Processor{}

	for _, trCfg := range trs {
		tr, err := newTransform(trCfg)
		if err != nil {
			return p, fmt.Errorf("in transform '%s': %w", trCfg.Name, err)
		}

		p.transforms = append(p.transforms, tr)
	}

	return p, nil
}

type transform struct {
	cfg pipeline.Transform

	condition *lang.Condition
	actions   []Action
}

func newTransform(cfg pipeline.Transform) (*transform, error) {
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
