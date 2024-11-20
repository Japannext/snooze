package transform

import (
	"context"

	"github.com/japannext/snooze/pkg/models"
)

type UnsetAction struct {
	Identity []string `yaml:""`
	Labels []string `yaml:"labels" json:"labels,omitempty"`
	AllIdentity bool `yaml:"all_identity" json:"allIdentity,omitempty"`
	AllLabels bool `yaml:"all_labels" json:"allLabels,omitempty"`
}

func (action *UnsetAction) Load() Transformation {
	return action
}

func (action *UnsetAction) Process(ctx context.Context, item *models.Log) (context.Context, error) {
	for _, key := range action.Labels {
		delete(item.Labels, key)
	}
	for _, key := range action.Identity {
		delete(item.Identity, key)
	}
	if action.AllIdentity {
		item.Identity = map[string]string{}
	}
	if action.AllLabels {
		item.Labels = map[string]string{}
	}
	return ctx, nil
}
