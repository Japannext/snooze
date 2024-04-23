package grouping

import (
	api "github.com/japannext/snooze/common/api/v2"
	"github.com/japannext/snooze/common/utils"
)

func Process(alert *api.Alert) error {
	var m map[string]string

	for _, rule := range computedRules {
		if rule.Condition.Test(alert) {
			for _, fi := range rule.Fields {
				v, found := fi.Get(alert)
				if found {
					m[fi.String()] = v
				}
			}
			break
		}
	}

	alert.GroupHash = utils.ComputeHash(m)
	return nil
}
