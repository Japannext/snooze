package silence

import (
	"github.com/japannext/snooze/common/schedule"
	"github.com/sirupsen/logrus"

	"github.com/japannext/snooze/common/lang"
)

type Rule struct {
	Name     string             `yaml:"name"`
	If       string             `yaml:"if"`
	Schedule *schedule.Schedule `yaml:",inline"`
}

type computedRule struct {
	Name     string
	Condition *lang.Condition
	Schedule schedule.Interface
}

func (r *computedRule) String() string {
	return r.Name
}

var computedRules []*computedRule

func compute(rule *Rule) *computedRule {
	condition, err := lang.NewCondition(rule.If)
	if err != nil {
		log.Fatal(err)
	}
	var s schedule.Interface
	if rule.Schedule == nil {
		s = &schedule.AlwaysSchedule{}
	} else {
		s, err = rule.Schedule.Resolve()
		if err != nil {
			log.Fatal(err)
		}
	}
	name := rule.Name
	if rule.Name == "" {
		name = rule.If
	}

	return &computedRule{
		Name: name,
		Condition: condition,
		Schedule: s,
	}
}

var log *logrus.Entry

func InitRules(rules []*Rule) {
	log = logrus.WithFields(logrus.Fields{"module": "silence"})
	for _, rule := range rules {
		computedRules = append(computedRules, compute(rule))
	}
}
