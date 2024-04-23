package transform

import (
	api "github.com/japannext/snooze/common/api/v2"
)

type ChildrenRule struct {
	Children []*Rule `yaml:"children"`
}

type computedChildren struct {
	children []*computedRule
}

func (rule *ChildrenRule) Compute() Interface {
	var items []*computedRule
	for _, child := range rule.Children {
		items = append(items, compute(child))
	}
	return &computedChildren{items}
}

func (rule *computedChildren) Process(alert *api.Alert) error {
	for _, child := range rule.children {
		if child.Condition.Test(alert) {
			if err := child.process.Process(alert); err != nil {
				return err
			}
		}
	}
	return nil
}
