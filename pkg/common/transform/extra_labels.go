package transform

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/lang"
)

type ExtraLabels struct {
	ExtraLabels map[string]string `yaml:"extra_labels"`

	internal struct {
		templates map[string]lang.Template
	}
}

func (tr *ExtraLabels) Load() error {
	var err error
	tr.internal.templates, err = lang.NewTemplateMap(tr.ExtraLabels)
	if err != nil {
		return err
	}
	return nil
}

func (tr *ExtraLabels) Transform(ctx context.Context, item *models.Log) error {
    for label, tpl := range tr.internal.templates {
        value, err := tpl.Execute(ctx, item)
        if err != nil {
            log.Warnf("failed to execute template `%s`", tr.ExtraLabels[label])
            return err
        }
        log.Debugf("Adding extra label: %s=%s", label, value)
        item.Labels[label] = value
    }
	return nil
}
