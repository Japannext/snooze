package transform

import (
	"context"

	"github.com/japannext/snooze/pkg/models"
)

type UnsetAction struct {
	Identity    []string `yaml:""`
	Labels      []string `json:"labels,omitempty"      yaml:"labels"`
	AllIdentity bool     `json:"allIdentity,omitempty" yaml:"all_identity"`
	AllLabels   bool     `json:"allLabels,omitempty"   yaml:"all_labels"`
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
