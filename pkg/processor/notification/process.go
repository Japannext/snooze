package notification

import (
	"context"
	set "github.com/deckarep/golang-set/v2"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/rabbitmq"
	"github.com/japannext/snooze/pkg/common/utils"
)

func Process(item *api.Log) error {

	ctx := context.Background()

	if item.Mute.SkipNotification {
		return nil
	}

	var queues = set.NewSet[*rabbitmq.NotificationQueue]()
	var merr = utils.NewMultiError("Failed to notify item trace_id=%s.")

	for _, rule := range computedRules {
		if rule.Condition != nil {
			v, err := rule.Condition.Match(ctx, item)
			if err != nil {
				return err
			}
			if v {
				for _, q := range rule.Queues {
					queues.Add(q)
				}
			}
		} else {
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
