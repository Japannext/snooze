package grouping

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/utils"
)

func Process(item *api.Log) error {
	ctx := context.Background()
	for _, rule := range computedRules {
		v, err := rule.Condition.Match(ctx, item)
		if err != nil {
			return err
		}
		if v {
			item.GroupLabels = make(map[string]string)
			for _, fi := range rule.Fields {
				v, err := fi.Get(ctx, item)
				if err != nil {
					logrus.Warnf("Failed to match %s: %s", fi, err)
					continue
				}
				item.GroupLabels[fi.String()] = v
			}
			item.GroupHash = utils.ComputeHash(item.GroupLabels)
			break
		}
	}

	return nil
}
