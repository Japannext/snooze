package processor

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/processor/grouping"
	"github.com/japannext/snooze/pkg/processor/notification"
	"github.com/japannext/snooze/pkg/processor/ratelimit"
	"github.com/japannext/snooze/pkg/processor/silence"
	"github.com/japannext/snooze/pkg/processor/store"
	"github.com/japannext/snooze/pkg/processor/transform"
)

type Processor struct{}

// For alert that will not be requeued, because their
// format is invalid, or they are poison messages.
type RejectedAlert struct {
	alert  *api.Alert
	reason string
}

func (r *RejectedAlert) Error() string {
	// return fmt.Sprintf("Rejected alert id=%s/%s reason=%s", r.alert.TraceID, r.alert.SpanID, r.reason)
	return fmt.Sprintf("Rejected alert: %s", r.reason)
}

func (p *Processor) Run() error {
	return ch.ConsumeForever(Process)
}

func (p *Processor) HandleStop() {
	ch.Cancel()
}

func Process(alert *api.Alert) error {
	log.Debugf("Start processing alert: %s", alert)
	if err := transform.Process(alert); err != nil {
		return err
	}
	if err := silence.Process(alert); err != nil {
		return err
	}
	if err := grouping.Process(alert); err != nil {
		return err
	}
	if err := ratelimit.Process(alert); err != nil {
		return err
	}
	if err := notification.Process(alert); err != nil {
		return err
	}
	if err := store.Process(alert); err != nil {
		return err
	}
	log.Debugf("End processing alert: %s", alert)
	return nil
}
