package profile

import (
	"fmt"
	"strings"
	"github.com/sirupsen/logrus"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

type Kv struct {
	Key string `yaml:"key"`
	Value string `yaml:"value"`
}

func FindValue(key string, item *api.Log) (value string, ok bool) {
	split := strings.Split(key, ".")
	if len(split) > 1 {
		switch (split[0]) {
			case "identity":
				value, ok = item.Identity[split[1]]
				break
			case "labels":
				value, ok = item.Labels[split[1]]
				break
			default:
				ok = false
		}
		return
	}
	value, ok = item.Labels[split[0]]
	return
}

type Rule struct {
	// Name of the profile group
	Name string `yaml:"name"`
	// The main condition for a log to match this rule. Used to
	// reduce the amount of processing by the use of maps.
	// Examples: process=sshd, service.name=keycloak, k8s.statefulset.name=postgresql
	Switch Kv `yaml:"switch"`
	// Patterns and actions to apply to logs matching this pattern
	Patterns []*Pattern `yaml:"patterns"`
}

func (rule *Rule) Startup() error {
	log.Debugf("[Startup] Profile %s", rule.Name)
	for _, pattern := range rule.Patterns {
		if err := pattern.Startup(); err != nil {
			return fmt.Errorf("in pattern %s: %w", pattern.Name, err)
		}
	}
	return nil
}

func (rule *Rule) Process(item *api.Log) bool {
	for _, pattern := range rule.Patterns {
		match, reject := pattern.Process(item)
		if reject {
			return true
		}
		if match {
			item.Profile = rule.Name
			item.Pattern = pattern.Name
			return false
		}
	}
	return false
}

var fastMapper *FastMapper

var log *logrus.Entry

func InitRules(rules []*Rule) {
	log = logrus.WithFields(logrus.Fields{"module": "profile"})
	fastMapper = NewFastMapper(rules)
}
