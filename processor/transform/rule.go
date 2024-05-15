package transform

import (
	"github.com/sirupsen/logrus"

	api "github.com/japannext/snooze/common/api/v2"
	"github.com/japannext/snooze/common/condition"
	"github.com/japannext/snooze/common/parser"
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
	Condition condition.Interface
	process   Interface
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

	c, err := parser.ParseCondition(rule.If)
	if err != nil {
		log.Fatal(err)
	}

	r := rule.Resolve()
	if r != nil {
		return &computedRule{c.Resolve(), r.Compute()}
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
