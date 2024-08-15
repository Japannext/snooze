package transform

import (
	"context"
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func Process(item *api.Log) error {
	ctx := context.Background()
	for _, rule := range computedRules {
		v, err := rule.Matcher.EvalBool(ctx, item)
		if err != nil {
			return err
		}
		if v {
			if err := rule.process.Process(item); err != nil {
				return err
			}
		}
	}

	return nil
}
