package transform

import (
	"context"
	api "github.com/japannext/snooze/pkg/common/api/v2"
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

func (rule *computedChildren) Process(item *api.Log) error {
	ctx := context.Background()
	for _, child := range rule.children {
		v, err := child.Matcher.EvalBool(ctx, item)
		if err != nil {
			return err
		}
		if v {
			if err := child.process.Process(item); err != nil {
				return err
			}
		}
	}
	return nil
}
