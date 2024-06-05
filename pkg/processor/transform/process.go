package transform

import (
	"context"
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func Process(alert *api.Alert) error {
	ctx := context.Background()
	for _, rule := range computedRules {
		v, err := rule.Matcher.EvalBool(ctx, alert)
		if err != nil {
			return err
		}
		if v {
			if err := rule.process.Process(alert); err != nil {
				return err
			}
		}
	}

	return nil
}
