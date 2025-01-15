package notification

import (
	"context"
	"fmt"

	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/tracing"
	"github.com/japannext/snooze/pkg/common/utils"
	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/processor/decision"
	"go.opentelemetry.io/otel"
	log "github.com/sirupsen/logrus"
)

type Processor struct {
	notifyQ             *mq.Pub
	notifications       []Notification
	defaultDestinations []models.Destination
}

type Config struct {
	// The notifications to execute. All that match will be executed.
	Notifications        []Notification `json:"notifications"        yaml:"notifications"`
	// The notifications to execute if nothing matches. All will be executed.
	DefaultDestinations []models.Destination `json:"defaultDestinations" yaml:"default_destinations"`
}

func New(cfg Config) (*Processor, error) {
	p := &Processor{}

	duplicate := utils.NewDuplicateChecker()
	notifyMap := map[string]Notification{}

	for index, notif := range cfg.Notifications {
		if notif.Name == "" {
			return p, fmt.Errorf("no name for notification at index #%d", index+1)
		}

		if err := duplicate.Check(notif.Name); err != nil {
			return p, fmt.Errorf("duplicate name '%s': %w", notif.Name, err)
		}

		if err := notif.Load(); err != nil {
			return p, fmt.Errorf("failed to load notification '%s': %w", notif.Name, err)
		}
		notifyMap[notif.Name] = notif
	}

	return p, nil
}

type Notification struct {
        Name         string               `json:"name"         yaml:"name"`
        If           string               `json:"if,omitempty" yaml:"if"`
        Destinations []models.Destination `json:"destinations" yaml:"destinations"`

        internal     struct {
                condition *lang.Condition
        }
}

func (notif *Notification) Load() error {
	if notif.If != "" {
		condition, err := lang.NewCondition(notif.If)
		if err != nil {
			return fmt.Errorf("failed to load condition `%s`: %s", notif.If, err)
		}

		notif.internal.condition = condition
	}

	return nil
}

func (p *Processor) getMatchingDestinations(ctx context.Context, item *models.Log) []models.Destination {
	destinations := utils.NewOrderedSet[models.Destination]()

	for _, notif := range p.notifications {
		if notif.internal.condition != nil {
			match, err := notif.internal.condition.MatchLog(ctx, item)
			if err != nil {
				continue
			}

			if !match {
				continue
			}

			log.Debugf("Matched notif '%s', destination(s): %s", notif.Name, notif.Destinations)
		}

		for _, dest := range notif.Destinations {
			destinations.Append(dest)
		}
	}

	return destinations.Items()
}

func (p *Processor) Process(ctx context.Context, item *models.Log) *decision.Decision {
	ctx, span := otel.Tracer("snooze").Start(ctx, "notification")
	defer span.End()

	if item.Status.SkipNotification {
		tracing.SetString(span, "notification.decision", "skipNotification")

		return decision.OK()
	}

	destinations := p.getMatchingDestinations(ctx, item)

	// Default destinations if nothing matches
	if len(destinations) == 0 {
		destinations = p.defaultDestinations
	}

	// Send notifications to destinations
	for _, dest := range destinations {
		notification := &models.Notification{
			Type:        "log",
			Destination: dest,
			Identity:    item.Identity,
			ItemID:      item.ID,
			Message:     item.Message,
		}

		// Send to notification queue
		subject := "NOTIFY." + dest.Queue
		if err := p.notifyQ.WithSubject(subject).Publish(ctx, notification); err != nil {
			log.Errorf("failed to queue notification to '%s:%s'", dest.Queue, dest.Profile)
			tracing.Error(span, err)

			return decision.Retry(err)
		}
	}

	return decision.OK()
}
