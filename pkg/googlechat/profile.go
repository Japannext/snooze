package googlechat

import (
	"os"

	log "github.com/sirupsen/logrus"
	chat "google.golang.org/api/chat/v1"
	"gopkg.in/yaml.v3"

	"github.com/japannext/snooze/pkg/googlechat/formatter"
	"github.com/japannext/snooze/pkg/models"
)

type Profiles struct {
	Profiles []*Profile `yaml:"profiles"`
}

type Profile struct {
	Name     string `yaml:"name"`
	Space    string `yaml:"space"`
	Timezone string `yaml:"timezone"`

	Format struct {
		Kind            string                     `yaml:"kind"`
		TemplateOptions *formatter.TemplateOptions `yaml:"-,squash"`
	} `yaml:"format"`

	Batch bool `yaml:"batch"`

	internal struct {
		formatter formatter.Interface
	}
}

func (p *Profile) Load() {
	switch p.Format.Kind {
	case "v1":
		p.internal.formatter = formatter.NewV1()
	case "template":
		p.internal.formatter = formatter.NewTemplate(p.Format.TemplateOptions)
	case "card":
		p.internal.formatter = formatter.NewCard()
	default:
		p.internal.formatter = formatter.NewV1()
	}
}

func (p *Profile) FormatToMessage(item *models.Notification) *chat.Message {
	msg, err := p.internal.formatter.Format(item)
	if err != nil {
		return formatter.FormatWithoutFail(item)
	}
	return msg
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

	for _, profile := range profileConfig.Profiles {
		profile.Load()
		if _, ok := profiles[profile.Name]; ok {
			log.Fatalf("Duplicate profile name '%s'", profile.Name)
		}
		profiles[profile.Name] = profile
	}
}
