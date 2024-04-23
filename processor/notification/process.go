package notification

import (
	"context"
	set "github.com/deckarep/golang-set/v2"

	api "github.com/japannext/snooze/common/api/v2"
	"github.com/japannext/snooze/common/rabbitmq"
	"github.com/japannext/snooze/common/utils"
)

func Process(alert *api.Alert) error {

	ctx := context.Background()

	if alert.Mute.Enabled {
		return nil
	}

	var queues = set.NewSet[*rabbitmq.NotificationQueue]()
	var merr = utils.NewMultiError("Failed to notify alert trace_id=%s.")

	for _, rule := range computedRules {
		if rule.Condition.Test(alert) {
			for _, q := range rule.Queues {
				queues.Add(q)
			}
		}
	}

	for q := range queues.Iter() {
		notif := &api.Notification{}
		if err := q.Publish(ctx, notif); err != nil {
			merr.AppendErr(err)
			continue
		}
	}

	if merr.HasErrors() {
		return merr
	}

	return nil
}
