package notification

import (
	"context"
	"time"
	"fmt"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/utils"
	"github.com/japannext/snooze/pkg/processor/tracing"
)

func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracing.TRACER.Start(ctx, "notification")
	defer span.End()

	if item.Mute.SkipNotification {
		return nil
	}

	// A set is necessary to avoid sending duplicates when 2 rules match
	// the same destination.
	var destinations = utils.NewOrderedSet[models.Destination]()
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
		notification := &models.Notification{
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
			mq.PublishAsync(fmt.Sprintf("NOTIFY.%s", dest.Queue), notification)
			if err := producer.Publish(notification); err != nil {
				merr.AppendErr(err)
				continue
			}
		}
		if _, err := opensearch.Store(ctx, models.NOTIFICATION_INDEX, notification); err != nil {
			merr.AppendErr(err)
			continue
		}
	}

	if merr.HasErrors() {
		return merr
	}

	return nil
}
