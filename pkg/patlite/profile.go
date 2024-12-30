package patlite

import (
	"os"

	pl "github.com/japannext/snooze/pkg/common/patlite"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var profiles = make(map[string]*Profile)

type Profiles struct {
	Profiles []*Profile `yaml:"profiles"`
}

/*
	Example configuration:

```yaml
---
profiles:
  - name: prod-patlite-red
    address: 10.x.x.x
    port: 10000
    lights:
    blue: "on"
    white: "on"
    red: "on"
    sound: "short"

```.
*/
type Profile struct {
	Name    string   `yaml:"name"`
	Address string   `yaml:"address"`
	Port    int      `yaml:"port"`
	State   pl.State `yaml:",inline"`
}

func loadProfiles() {
	data, err := os.ReadFile(env.ProfilePath)
	if err != nil {
		log.Fatal(err)
	}
	var profileConfig Profiles
	if err := yaml.Unmarshal(data, &profileConfig); err != nil {
		log.Fatal(err)
	}

	for _, profile := range profileConfig.Profiles {
		if _, ok := profiles[profile.Name]; ok {
			log.Fatalf("Duplicate profile name '%s'", profile.Name)
		}
		profiles[profile.Name] = profile
	}
}
