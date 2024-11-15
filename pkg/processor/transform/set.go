package transform

import (
	"context"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/lang"
)

type SetAction struct {
	Identity map[string]string `yaml:"identity" json:"identity,omitempty"`
	Labels map[string]string `yaml:"labels" json:"labels,omitempty"`
	Message string `yaml:"message" json:"message,omitempty"`

	internal struct {
		identity map[string]*lang.Template
		labels map[string]*lang.Template
		message *lang.Template
	}
}

func (action *SetAction) Load() Transformation {
	var err error
	if len(action.Labels) > 0 {
		if action.internal.labels, err = lang.NewTemplateMap(action.Labels); err != nil {
			log.Errorf("bad template: %s", err)
		}
	}
	if len(action.Identity) > 0 {
		if action.internal.identity, err = lang.NewTemplateMap(action.Identity); err != nil {
			log.Errorf("bad template: %s", err)
		}
	}
	if len(action.Message) > 0 {
		if action.internal.message, err = lang.NewTemplate(action.Message); err != nil {
			log.Errorf("bad template: %s", err)
		}
	}
	return action
}

func (action *SetAction) Process(ctx context.Context, item *models.Log) (context.Context, error) {
	for key, tpl := range action.internal.labels {
		value, err := tpl.Execute(ctx, item)
		if err != nil {
			log.Warnf("failed to execute template `%s`: %s", action.Labels[key], err)
			continue
		}
		item.Labels[key] = value
	}
	for key, tpl := range action.internal.identity {
		value, err := tpl.Execute(ctx, item)
		if err != nil {
			log.Warnf("failed to execute template `%s`: %s", action.Identity[key], err)
			continue
		}
		item.Identity[key] = value
	}
	if tpl := action.internal.message; tpl != nil {
		value, err := tpl.Execute(ctx, item)
		if err != nil {
			log.Warnf("failed to execute template `%s`: %s", action.Message, err)
		}
		item.Message = value
	}
	return ctx, nil
}
