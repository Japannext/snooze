package notification

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"

	"github.com/japannext/snooze/components/processor/decision"
	"github.com/japannext/snooze/internal/lang"
	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/mq"
	"github.com/japannext/snooze/internal/orderedset"
	"github.com/japannext/snooze/internal/pipeline"
	"github.com/japannext/snooze/internal/tracing"
	"github.com/japannext/snooze/internal/utils"
)

type Processor struct {
	notifyQ             *mq.Pub
	notifications       []Notification
	defaultDestinations []models.Destination
}

func New(cfgs []pipeline.Notification, dests []models.Destination) (*Processor, error) {
	p := &Processor{}

	duplicate := utils.NewDuplicateChecker()
	notifyMap := map[string]Notification{}

	for index, cfg := range cfgs {
		if cfg.Name == "" {
			return &Processor{}, fmt.Errorf("no name for notification at index #%d", index+1)
		}

		if err := duplicate.Check(cfg.Name); err != nil {
			return &Processor{}, fmt.Errorf("duplicate name '%s': %w", cfg.Name, err)
		}

		notif, err := NewNotification(cfg)
		if err != nil {
			return &Processor{}, fmt.Errorf("failed to load notification '%s': %w", cfg.Name, err)
		}
		notifyMap[cfg.Name] = notif
	}

	return p, nil
}

type Notification struct {
	pipeline.Notification
	condition *lang.Condition
}

func NewNotification(cfg pipeline.Notification) (Notification, error) {
	notif := Notification{Notification: cfg}
	if cfg.If != "" {
		condition, err := lang.NewCondition(cfg.If)
		if err != nil {
			return Notification{}, fmt.Errorf("failed to load condition `%s`: %s", cfg.If, err)
		}

		notif.condition = condition
	}

	return notif, nil
}

func (p *Processor) getMatchingDestinations(ctx context.Context, item *models.Log) []models.Destination {
	destinations := orderedset.New[models.Destination]()

	for _, notif := range p.notifications {
		if notif.condition != nil {
			match, err := notif.condition.MatchLog(ctx, item)
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
