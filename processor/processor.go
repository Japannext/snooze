package processor

import (
  amqp "github.com/rabbitmq/amqp091-go"

  "github.com/japannext/snooze/common/rabbitmq"
  "github.com/japannext/snooze/common/api/v2"
)

type Processor struct {}

// For alert that will not be requeued, because their
// format is invalid, or they are poison messages.
type RejectedAlert struct {
  alert *v2.Alert
  reason string
}
func (r *RejectedAlert) Error() string {
  return fmt.Sprintf("Rejected alert id=%s/%s reason=%s", r.alert.TraceID, r.alert.SpanID, r.reason)
}

func (p *Processor) Run() error {
  dd, err := pq.Consume()
  if err != nil {
    return err
  }
  for _, d := range dd {
    var alert *v2.Alert
    if err := json.UnMarshal(d.Body, &alert); err != nil {
      log.Warn("Rejected malformed alert body:\n%.20s\n[...]", d.Body)
      d.Reject(false)
    }
    pipeline := pipelines[pname]
    newItem, err := pipeline.Process(item)
    if errors.As(err, &RejectedAlert{}) {
      log.Error(err)
      d.Reject(false)
    }
    if err != nil {
      log.Warn("Failed to process alert id=%s/%s, err=%s", alert.TraceID, alert.SpanID, err)
      d.Reject(true)
    }
    d.Ack()
  }
  return nil
}

func (p *Processor) HandleStop() {
  pq.Cancel()
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
  if err := save.Process(alert); err != nil {
    return err
  }
}
