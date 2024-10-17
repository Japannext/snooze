package notification

import (
	"context"
	"time"
	"fmt"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/utils"
)

func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracer.Start(ctx, "notification")
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
		subject := fmt.Sprintf("NOTIFY.%s", dest.Queue)
		if dest.Queue != "dummy" {
			if err := notifyQ.WithSubject(subject).Publish(ctx, notification); err != nil {
				log.Warnf("failed to notify: %s", err)
			}
		}
	}

	if merr.HasErrors() {
		return merr
	}

	return nil
}
