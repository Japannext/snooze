package profile

import (
	"context"

	"github.com/japannext/snooze/pkg/models"
)

type Profile struct {
	// Name of the profile group
	Name string `json:"name" yaml:"name"`
	// The main condition for a log to match this rule. Used to
	// reduce the amount of processing by the use of maps.
	// Examples: process=sshd, service.name=keycloak, k8s.statefulset.name=postgresql
	Switch Kv `json:"switch" yaml:"switch"`
	// Patterns and actions to apply to logs matching this pattern
	Patterns []*Pattern `json:"patterns" yaml:"patterns"`
}

func (prf *Profile) Load() {
	for _, pattern := range prf.Patterns {
		if err := pattern.Load(); err != nil {
			log.Fatalf("in profile=%s, pattern=%s: %s", prf.Name, pattern.Name, err)
		}
	}
}

func (prf *Profile) Process(ctx context.Context, item *models.Log) bool {
	for _, pattern := range prf.Patterns {
		match, reject := pattern.Process(ctx, item)
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
