package processor

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/processor/grouping"
	"github.com/japannext/snooze/pkg/processor/profile"
	"github.com/japannext/snooze/pkg/processor/notification"
	"github.com/japannext/snooze/pkg/processor/ratelimit"
	"github.com/japannext/snooze/pkg/processor/silence"
	"github.com/japannext/snooze/pkg/processor/store"
	"github.com/japannext/snooze/pkg/processor/transform"
)

type Processor struct{}

// For item that will not be requeued, because their
// format is invalid, or they are poison messages.
type RejectedLog struct {
	item  *api.Log
	reason string
}

func (r *RejectedLog) Error() string {
	// return fmt.Sprintf("Rejected item id=%s/%s reason=%s", r.item.TraceID, r.item.SpanID, r.reason)
	return fmt.Sprintf("Rejected item: %s", r.reason)
}

func (p *Processor) Run() error {
	return ch.ConsumeForever(Process)
}

func (p *Processor) HandleStop() {
	ch.Stop()
}

func Process(item *api.Log) error {
	log.Debugf("Start processing item: %s", item)
	if err := transform.Process(item); err != nil {
		return err
	}
	if err := silence.Process(item); err != nil {
		return err
	}
	if err := profile.Process(item); err != nil {
		return err
	}
	if err := grouping.Process(item); err != nil {
		return err
	}
	if err := ratelimit.Process(item); err != nil {
		return err
	}
	if err := notification.Process(item); err != nil {
		return err
	}
	if err := store.Process(item); err != nil {
		return err
	}
	log.Debugf("End processing item: %s", item)
	return nil
}
