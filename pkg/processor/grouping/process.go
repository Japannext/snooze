package grouping

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/utils"
)

func Process(alert *api.Alert) error {
	ctx := context.Background()
	for _, rule := range computedRules {
		v, err := rule.Condition.Match(ctx, alert)
		if err != nil {
			return err
		}
		if v {
			alert.GroupLabels = make(map[string]string)
			for _, fi := range rule.Fields {
				v, err := fi.Get(ctx, alert)
				if err != nil {
					logrus.Warnf("Failed to match %s: %s", fi, err)
					continue
				}
				alert.GroupLabels[fi.String()] = v
			}
			alert.GroupHash = utils.ComputeHash(alert.GroupLabels)
			break
		}
	}

	return nil
}
