package lang

import (
	"bytes"
	"context"
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

func NewTemplateMap(rawMap map[string]string) (map[string]*Template, error) {
	results := make(map[string]*Template)
	for key, text := range rawMap {
		tpl, err := NewTemplate(text)
		if err != nil {
			return nil, fmt.Errorf("invalid template at key='%s' template='%s': %w", key, text, err)
		}
		results[key] = tpl
	}
	return results, nil
}

// Extract the capture from the context
func getCapture(ctx context.Context) map[string]string {
	capture, ok := ctx.Value("capture").(map[string]string)
	if !ok {
		return map[string]string{}
	}
	return capture
}

func getMappings(ctx context.Context) map[string]string {
	mappings, ok := ctx.Value("mappings").(map[string]string)
	if !ok {
		return map[string]string{}
	}
	return mappings
}

func (t *Template) Execute(ctx context.Context, item models.HasContext) (string, error) {
	data := item.Context()
	data["capture"] = getCapture(ctx)
	data["mappings"] = getMappings(ctx)

	var buf bytes.Buffer

	err := t.template.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
