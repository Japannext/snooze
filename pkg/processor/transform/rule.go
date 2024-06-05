package transform

import (
	"github.com/PaesslerAG/gval"
	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/sirupsen/logrus"
)

type Rule struct {
	If       string        `yaml:"if"`
	Children *ChildrenRule `yaml:",inline"`
	Set      *SetRule      `yaml:"set,omitempty"`
}

type Interface interface {
	Process(*api.Alert) error
}

type Computable interface {
	Compute() Interface
}

type computedRule struct {
	Matcher gval.Evaluable
	process Interface
}

var computedRules []*computedRule

func (rule *Rule) Resolve() Computable {
	if rule.Children != nil {
		return rule.Children
	}
	if rule.Set != nil {
		return rule.Set
	}
	return nil
}

func compute(rule *Rule) *computedRule {

	matcher, err := gval.Full().NewEvaluable(rule.If)
	if err != nil {
		log.Fatal(err)
	}

	r := rule.Resolve()
	if r != nil {
		return &computedRule{matcher, r.Compute()}
	}
	return nil
}

var log *logrus.Entry

func InitRules(rules []*Rule) {
	log = logrus.WithFields(logrus.Fields{"module": "transform"})
	for _, rule := range rules {
		computedRules = append(computedRules, compute(rule))
	}
}
