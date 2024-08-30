package profile

import (
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

type Profile struct {
	// Name of the profile group
	Name string `yaml:"name"`
	// The main condition for a log to match this rule. Used to
	// reduce the amount of processing by the use of maps.
	// Examples: process=sshd, service.name=keycloak, k8s.statefulset.name=postgresql
	Switch Kv `yaml:"switch"`
	// Patterns and actions to apply to logs matching this pattern
	Patterns []*Pattern `yaml:"patterns"`
}

func (prf *Profile) Load() {
	for _, pattern := range prf.Patterns {
		if err := pattern.Load(); err != nil {
			log.Fatalf("in profile=%s, pattern=%s: %s", prf.Name, pattern.Name, err)
		}
	}
}

func (prf *Profile) Process(item *api.Log) bool {
	for _, pattern := range prf.Patterns {
		match, reject := pattern.Process(item)
		if reject {
			return true
		}
		if match {
			item.Profile = prf.Name
			item.Pattern = pattern.Name
			return false
		}
	}
	return false
}
