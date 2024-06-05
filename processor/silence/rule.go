package silence

import (
	"github.com/PaesslerAG/gval"
	"github.com/japannext/snooze/common/schedule"
	"github.com/sirupsen/logrus"
)

type Rule struct {
	Name     string             `yaml:"name"`
	If       string             `yaml:"if"`
	Schedule *schedule.Schedule `yaml:",inline"`
}

type computedRule struct {
	Name     string
	Matcher  gval.Evaluable
	Schedule schedule.Interface
}

func (r *computedRule) String() string {
	return r.Name
}

var computedRules []*computedRule

func compute(rule *Rule) *computedRule {
	matcher, err := gval.Full().NewEvaluable(rule.If)
	if err != nil {
		log.Fatal(err)
	}
	s, err := rule.Schedule.Resolve()
	if err != nil {
		log.Fatal(err)
	}
	name := rule.Name
	if rule.Name == "" {
		name = rule.If
	}

	return &computedRule{
		Name:     name,
		Matcher:  matcher,
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
