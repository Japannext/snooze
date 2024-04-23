package processor

import (
  "fmt"

  log "github.com/sirupsen/logrus"

  "github.com/japannext/snooze/processor/transform"
  "github.com/japannext/snooze/processor/silence"
  "github.com/japannext/snooze/processor/grouping"
  "github.com/japannext/snooze/processor/ratelimit"
  "github.com/japannext/snooze/processor/notification"
  "github.com/japannext/snooze/processor/store"
  api "github.com/japannext/snooze/common/api/v2"
)

type Processor struct {}

// For alert that will not be requeued, because their
// format is invalid, or they are poison messages.
type RejectedAlert struct {
  alert *api.Alert
  reason string
}
func (r *RejectedAlert) Error() string {
  // return fmt.Sprintf("Rejected alert id=%s/%s reason=%s", r.alert.TraceID, r.alert.SpanID, r.reason)
  return fmt.Sprintf("Rejected alert: %s", r.reason)
}

func (p *Processor) Run() error {
  msgs, err := ch.Consume()
  if err != nil {
    return err
  }
  for msg := range msgs {
    alert := msg.Alert
    d := msg.Delivery
    err = Process(alert)
    if err != nil {
      log.Error(err)
      d.Reject(false)
    }
    d.Ack(false)
  }
  return nil
}

func (p *Processor) HandleStop() {
  ch.Cancel()
}

func Process(alert *api.Alert) error {
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
  return nil
}
