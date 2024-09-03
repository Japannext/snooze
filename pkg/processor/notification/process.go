package notification

import (
	"context"
	"time"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/utils"
)

func Process(ctx context.Context, item *api.Log) error {

	if item.Mute.SkipNotification {
		return nil
	}

	// A set is necessary to avoid sending duplicates when 2 rules match
	// the same destination.
	var destinations = utils.NewOrderedSet[api.Destination]()
	var merr = utils.NewMultiError("Failed to notify item trace_id=%s.")

	for _, notif := range notifications {
		if notif.internal.condition != nil {
			match, err := notif.internal.condition.MatchLog(ctx, item)
			if err != nil {
				return err
			}
			if !match {
				continue
			}
			log.Debugf("Matched notif '%s', destination(s): %s", notif.If, notif.Destinations)
		}
		destinations.AppendMany(notif.Destinations)
	}
	// Defaults
	if len(destinations.Items()) == 0 {
		log.Debugf("Match no notif, default destinations: %s", defaultDestinations)
		destinations.AppendMany(defaultDestinations)
	}

	for _, dest := range destinations.Items() {
		log.Debugf("sending to destination `%s`", dest)
		notification := &api.Notification{
			TimestampMillis: uint64(time.Now().UnixMilli()),
			Destination: dest,
			LogUID: item.ID,
			Body: map[string]string{
				"message": item.Message,
			},
		}
		if dest.Queue != "dummy" {
			producer, found := producers[dest.Queue]
			if !found {
				log.Errorf("Producer for queue '%s' not found! This should not happen!", dest.Queue)
				continue
			}
			if err := producer.Publish(notification); err != nil {
				merr.AppendErr(err)
				continue
			}
		}
		if _, err := opensearch.StoreNotification(ctx, notification); err != nil {
			merr.AppendErr(err)
			continue
		}
	}

	if merr.HasErrors() {
		return merr
	}

	return nil
}
