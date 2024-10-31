package formatter

import (
	"bytes"
	"text/template"

	chat "google.golang.org/api/chat/v1"
	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/models"
)

type TemplateOptions struct {
	Template string `yaml:"template"`
}

type Template struct {
	template *template.Template
}

func NewTemplate(opts *TemplateOptions) *Template {

	if opts.Template == "" {
		log.Fatalf("empty `template` provided")
	}

	tpl, err := template.New("template").Parse(opts.Template)
	if err != nil {
		log.Fatalf("invalid template `%s`: %s", opts.Template, err)
	}
	return &Template{
		template: tpl,
	}
}

func (tpl *Template) Format(item *models.Notification) (*chat.Message, error) {
    var buf bytes.Buffer
    err := tpl.template.Execute(&buf, item.Context())
    if err != nil {
		return &chat.Message{}, err
    }

	return &chat.Message{Text: buf.String()}, nil
}
