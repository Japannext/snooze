package notification

import (
	"context"
	"fmt"

	"github.com/japannext/snooze/pkg/common/tracing"
	"github.com/japannext/snooze/pkg/models"
)

func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracer.Start(ctx, "notification")
	defer span.End()

	if item.Status.SkipNotification {
		tracing.SetString(span, "notification.decision", "skipNotification")
		return nil
	}

	// A set is necessary to avoid sending duplicates when 2 rules match
	// the same destination.
	// var destinations = utils.NewOrderedSet[models.Destination]()
	destinations := make(map[models.Destination]bool)

	tracing.SetInt(span, "notifications.number", len(notifications))

	for _, notif := range notifications {
		if notif.internal.condition != nil {
			match, err := notif.internal.condition.MatchLog(ctx, item)
			if err != nil {
				tracing.Error(span, err)
				return err
			}
			if !match {
				tracing.SetBool(span, fmt.Sprintf("notification.%s.match", notif.Name), false)
				continue
			}
			log.Debugf("Matched notif '%s', destination(s): %s", notif.Name, notif.Destinations)
		}
		tracing.SetBool(span, fmt.Sprintf("notification.%s.match", notif.Name), true)

		for _, dest := range notif.Destinations {
			if dest.Queue == "dummy" {
				continue
			}
			// Mechanism to avoid send 2 notifications to the same destination
			// if it matches multiple times
			if _, found := destinations[dest]; found {
				continue
			}
			destinations[dest] = true

			notification := &models.Notification{
				Type:        "log",
				Destination: dest,
				Identity:    item.Identity,
				ItemID:      item.ID,
				Message:     item.Message,
			}

			// Send to notification queue
			subject := fmt.Sprintf("NOTIFY.%s", dest.Queue)
			if err := notifyQ.WithSubject(subject).Publish(ctx, notification); err != nil {
				log.Warnf("failed to queue notification to '%s:%s'", dest.Queue, dest.Profile)
				tracing.SetString(span, fmt.Sprintf("notification.%s.%s:%s", notif.Name, dest.Queue, dest.Profile), "failed")
				tracing.Error(span, err)
			}
			tracing.SetString(span, fmt.Sprintf("notification.%s.%s:%s", notif.Name, dest.Queue, dest.Profile), "sent")
		}
	}

	return nil
}
