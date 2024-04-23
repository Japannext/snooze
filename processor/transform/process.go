package transform

import (
	api "github.com/japannext/snooze/common/api/v2"
)

func Process(alert *api.Alert) error {
	for _, rule := range computedRules {
		if rule.Condition.Test(alert) {
			if err := rule.process.Process(alert); err != nil {
				return err
			}
		}
	}

	return nil
}
