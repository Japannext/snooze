package unset

import (
	"context"

	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/pipeline"
)

type Action struct {
	cfg *pipeline.Unset
}

func New(cfg *pipeline.Unset) (*Action, error) {
	return &Action{cfg: cfg}, nil
}

func (action *Action) Process(ctx context.Context, item *models.Log) (context.Context, error) {
	for _, key := range action.cfg.Labels {
		delete(item.Labels, key)
	}

	for _, key := range action.cfg.Identity {
		delete(item.Identity, key)
	}

	if action.cfg.AllIdentity {
		item.Identity = map[string]string{}
	}

	if action.cfg.AllLabels {
		item.Labels = map[string]string{}
	}

	return ctx, nil
}
