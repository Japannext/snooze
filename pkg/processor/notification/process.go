package notification

import (
	"context"
	"time"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/utils"
)

func Process(item *api.Log) error {

	ctx := context.Background()

	if item.Mute.SkipNotification {
		return nil
	}

	var queues = utils.NewOrderedStringSet()
	var merr = utils.NewMultiError("Failed to notify item trace_id=%s.")

	for _, rule := range computedRules {
		if rule.internal.condition != nil {
			match, err := rule.internal.condition.MatchLog(ctx, item)
			if err != nil {
				return err
			}
			if !match {
				continue
			}
		}
		queues.AppendMany(rule.Channels)
	}

	for _, queue := range queues.Items() {
		log.Debugf("sending to queue `%s`", queue)
		notification := &api.Notification{
			TimestampMillis: uint64(time.Now().UnixMilli()),
			Destination: api.Destination{Name: queue},
			LogUID: item.ID,
			Body: map[string]string{
				"message": item.Message,
			},
		}
		producer, found := producers[queue]
		if !found {
			log.Errorf("Producer for queue '%s' not found! This should not happen!", queue)
			continue
		}
		if err := producer.Publish(notification); err != nil {
			merr.AppendErr(err)
			continue
		}
		if err := opensearch.LogStore.StoreNotification(notification); err != nil {
			merr.AppendErr(err)
			continue
		}
	}

	if merr.HasErrors() {
		return merr
	}

	return nil
}
