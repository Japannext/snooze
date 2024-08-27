package mail

import (
	"fmt"
	"os"
	"text/template"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const DEFAULT_TEMPLATE_TEXT = `
Subject: Snooze notification ({{ .Destination.Profile }})

Your received a snooze notification!
Body = {{ .Body }}
Labels = {{ .Labels }}
DocumentationURL = {{ .DocumentationURL }}
SnoozeURL = {{ .SnoozeURL }}
`

var defaultTemplate *template.Template

var profiles = make(map[string]*Profile)

type ProfileConfig struct {
	Profiles []*Profile `yaml:"profiles"`
}

type Profile struct {
	Name string `yaml:"name"`
	From string `yaml:"from"`
	To string `yaml:"to"`
	Template string `yaml:"template"`
	Server string `yaml:"server"`
	Port int `yaml:"port"`

	internal struct {
		template *template.Template
		isDefault bool
	}
}

func (p *Profile) Startup() error {
	if p.From == "" {
		p.From = config.DefaultSender
	}
	if p.Server == "" {
		p.Server = config.Server
	}
	if p.Port == 0 {
		p.Port = config.Port
	}
	if p.Template != "" {
		var err error
		p.internal.template, err = template.New("template").Parse(p.Template)
		if err != nil {
			return fmt.Errorf("invalid template for profile '%s': %w", p.Name, err)
		}
	} else {
		p.internal.template = defaultTemplate
		p.internal.isDefault = true
	}

	return nil
}

func loadProfiles() {
	data, err := os.ReadFile(config.ProfilePath)
	if err != nil {
		log.Fatal(err)
	}
	var profileConfig ProfileConfig
	if err := yaml.Unmarshal(data, &profileConfig); err != nil {
		log.Fatal(err)
	}

	defaultTemplate, err = template.New("default_template").Parse(DEFAULT_TEMPLATE_TEXT)
	if err != nil {
		log.Fatalf("default template failed to be parsed: %s", err)
	}

	for _, profile := range profileConfig.Profiles {
		profile.Startup()
		if _, ok := profiles[profile.Name]; ok {
			log.Fatalf("Duplicate profile name '%s'", profile.Name)
		}
		profiles[profile.Name] = profile
	}

}
