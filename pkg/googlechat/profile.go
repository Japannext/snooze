package googlechat

import (
    "os"
    "text/template"

    log "github.com/sirupsen/logrus"
    "gopkg.in/yaml.v3"
)

const DEFAULT_TEMPLATE_TEXT = `
At @{{ .TimestampMillis }}:
{{ .Body }}
{{ if .DocumentationURL }}Doc: {{ .DocumentationURL }}{{ end }}
`

var defaultTemplate *template.Template

type Profiles struct {
	Profiles []*Profile `yaml:"profiles"`
}

type Profile struct {
    Name string `yaml:"name"`
	Space string `yaml:"space"`
	Template string `yaml:"template"`
	Batch bool `yaml:"batch"`

    internal struct {
		template *template.Template
		isDefault bool
    }
}

func (p *Profile) Load() {
	if p.Template != "" {
		var err error
		p.internal.template, err = template.New("template").Parse(p.Template)
		if err != nil {
			log.Fatalf("invalid template for profile '%s': %s", p.Name, err)
		}
	} else {
		p.internal.template = defaultTemplate
		p.internal.isDefault = true
	}
}

var profiles = map[string]*Profile{}

func loadProfiles() {
    data, err := os.ReadFile(config.ProfilePath)
    if err != nil {
        log.Fatal(err)
    }
    var profileConfig Profiles
    if err := yaml.Unmarshal(data, &profileConfig); err != nil {
        log.Fatal(err)
    }

    defaultTemplate, err = template.New("default_template").Parse(DEFAULT_TEMPLATE_TEXT)
    if err != nil {
        log.Fatalf("default template failed to be parsed: %s", err)
    }

    for _, profile := range profileConfig.Profiles {
		profile.Load()
        if _, ok := profiles[profile.Name]; ok {
            log.Fatalf("Duplicate profile name '%s'", profile.Name)
        }
        profiles[profile.Name] = profile
    }
}
