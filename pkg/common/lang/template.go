package lang

import (
	"bytes"

	"text/template"
)

type Template struct {
	template *template.Template
}

func NewTemplate(name, raw string) (*Template, error) {
	tmpl, err := template.New(name).Parse(raw)
	if err != nil {
		return nil, err
	}
	return &Template{tmpl}, nil
}

func (tmpl *Template) Execute(data any) (string, error) {
	var buf bytes.Buffer
	err := tmpl.template.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
