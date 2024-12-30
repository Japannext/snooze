package profile

import (
	"bytes"
	"fmt"

	"text/template"

	"github.com/japannext/snooze/pkg/models"
)

type Template struct {
	template *template.Template
}

func NewTemplate(raw string) (*Template, error) {
	t, err := template.New("log-pattern").Parse(raw)
	if err != nil {
		return nil, err
	}
	return &Template{t}, nil
}

func NewTemplateMap(rawMap map[string]string) (map[string]Template, error) {
	results := make(map[string]Template)
	for key, text := range rawMap {
		tpl, err := NewTemplate(text)
		if err != nil {
			return nil, fmt.Errorf("invalid template at key='%s' template='%s': %w", key, text, err)
		}
		results[key] = *tpl
	}
	return results, nil
}

func (t *Template) Execute(item *models.Log, capture map[string]string) (string, error) {
	data := map[string]interface{}{
		"actualTime": item.ActualTime,
		"labels":     item.Labels,
		"identity":   item.Identity,
		"message":    item.Message,
	}
	if len(capture) > 0 {
		data["capture"] = capture
	}
	var buf bytes.Buffer
	err := t.template.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
