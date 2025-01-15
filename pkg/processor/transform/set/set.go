package set

import (
	"context"
	"fmt"

	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/models"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Identity map[string]string `json:"identity,omitempty" yaml:"identity"`
	Labels   map[string]string `json:"labels,omitempty"   yaml:"labels"`
	Message  string            `json:"message,omitempty"  yaml:"message"`
}

type Action struct {
	cfg *Config

	identity map[string]*lang.Template
	labels   map[string]*lang.Template
	message  *lang.Template
}

func New(cfg *Config) (*Action, error) {
	action := &Action{
		cfg: cfg,
	}

	var err error
	if len(cfg.Labels) > 0 {
		action.labels, err = lang.NewTemplateMap(cfg.Labels)
		if err != nil {
			return action, fmt.Errorf("bad template: %w", err)
		}
	}

	if len(cfg.Identity) > 0 {
		action.identity, err = lang.NewTemplateMap(cfg.Identity)
		if err != nil {
			return action, fmt.Errorf("bad template: %s", err)
		}
	}

	if len(cfg.Message) > 0 {
		action.message, err = lang.NewTemplate(cfg.Message)
		if err != nil {
			return action, fmt.Errorf("bad template: %s", err)
		}
	}

	return action, nil
}

func (action *Action) Process(ctx context.Context, item *models.Log) (context.Context, error) {
	for key, tpl := range action.labels {
		value, err := tpl.Execute(ctx, item)
		if err != nil {
			log.Warnf("failed to execute template `%s`: %s", action.cfg.Labels[key], err)

			continue
		}

		item.Labels[key] = value
	}

	for key, tpl := range action.identity {
		value, err := tpl.Execute(ctx, item)
		if err != nil {
			log.Warnf("failed to execute template `%s`: %s", action.cfg.Identity[key], err)

			continue
		}

		item.Identity[key] = value
	}

	if tpl := action.message; tpl != nil {
		value, err := tpl.Execute(ctx, item)
		if err != nil {
			log.Warnf("failed to execute template `%s`: %s", action.cfg.Message, err)
		}

		item.Message = value
	}

	return ctx, nil
}
