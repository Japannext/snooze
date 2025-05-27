package transform

import (
	"context"
	"fmt"

	"github.com/japannext/snooze/components/processor/transform/regex"
	"github.com/japannext/snooze/components/processor/transform/set"
	"github.com/japannext/snooze/components/processor/transform/unset"
	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/pipeline"
)

type Action interface {
	Process(ctx context.Context, item *models.Log) (context.Context, error)
}

func NewAction(cfg pipeline.Action) (Action, error) {
	switch {
	case cfg.Set != nil:
		return set.New(cfg.Set)
	case cfg.Unset != nil:
		return unset.New(cfg.Unset)
	case cfg.Regex != nil:
		return regex.New(cfg.Regex)
	default:
		return nil, fmt.Errorf("empty action")
	}
}
