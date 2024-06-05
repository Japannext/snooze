package lang

import (
	"context"
	"github.com/PaesslerAG/gval"

	api "github.com/japannext/snooze/common/api/v2"
)

type Condition struct {
	raw  string
	gval gval.Evaluable
}

func NewCondition(raw string) (*Condition, error) {
	e, err := gval.Full().NewEvaluable(raw)
	if err != nil {
		return &Condition{}, err
	}
	return &Condition{raw, e}, nil
}

func (c *Condition) String() string {
	return c.raw
}

func (c *Condition) Match(ctx context.Context, alert *api.Alert) (bool, error) {
	return c.gval.EvalBool(ctx, map[string]interface{}{"alert": alert})
}
