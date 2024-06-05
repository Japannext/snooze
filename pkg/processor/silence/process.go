package silence

import (
	"context"
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func Process(alert *api.Alert) error {

	ctx := context.Background()

	for _, rule := range computedRules {
		v, err := rule.Condition.Match(ctx, alert)
		if err != nil {
			return err
		}
		if v {
			// Silence the alert
			alert.Mute.Enabled = true
			alert.Mute.Component = "silence"
			alert.Mute.Rule = rule.String()
			alert.Mute.SkipNotification = true
			break
		}
	}

	return nil
}
