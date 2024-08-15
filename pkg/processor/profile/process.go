package profile

import (
	// "github.com/sirupsen/logrus"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func Process(item *api.Log) error {
	for _, rule := range fastMapper.GetRules(item) {
		reject := rule.Process(item)
		if reject {
			return nil
		}
	}

	return nil
}
