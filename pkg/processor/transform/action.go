package transform

import (
	"context"
	"fmt"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/processor/transform/set"
	"github.com/japannext/snooze/pkg/processor/transform/unset"
	"github.com/japannext/snooze/pkg/processor/transform/regex"
)

type ActionConfig struct {
	Set   *set.Config   `json:"set,omitempty"   yaml:"set"`
	Unset *unset.Config `json:"unset,omitempty" yaml:"unset"`
	Regex *regex.Config `json:"regex,omitempty" yaml:"regex"`
}

type ActionInterface interface {
	Process(ctx context.Context, item *models.Log) (context.Context, error)
}

func NewAction(cfg *ActionConfig) (ActionInterface, error) {
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
