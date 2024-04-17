package transform

import (
  "fmt"

  api "github.com/japannext/snooze/common/api/v2"
  condition "github.com/japannext/snooze/common/condition"
  "github.com/japannext/snooze/common/ref"
)

type TransformConfig struct {
  Condition string
}

func (act *TransformAction) Init() error {
  return nil
}

func (a *TransformAction) Process(item api.Alert) error {
  return nil
}

type TransformInterface interface {
  Execute(*api.Alert)
  String() string
}

type SetOp struct {
  Ref ref.Reference
  Value string
}

func (op *SetOp) Execute(a *api.Alert) {
  op.Ref.Set(a, op.Value)
}

func (op *SetOp) String() string {
  return fmt.Sprintf("%s -> %s", op.Ref.String(), op.Value)
}
