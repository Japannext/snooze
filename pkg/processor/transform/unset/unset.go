package unset

import (
	"context"

	"github.com/japannext/snooze/pkg/models"
)

type Config struct {
	Identity    []string `yaml:""`
	Labels      []string `json:"labels,omitempty"      yaml:"labels"`
	AllIdentity bool     `json:"allIdentity,omitempty" yaml:"all_identity"`
	AllLabels   bool     `json:"allLabels,omitempty"   yaml:"all_labels"`
}

type Action struct {
	cfg *Config
}

func New(cfg *Config) (*Action, error) {
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
