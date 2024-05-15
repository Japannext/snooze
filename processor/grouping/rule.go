package grouping

import (
	"github.com/japannext/snooze/common/condition"
	"github.com/japannext/snooze/common/field"
	"github.com/japannext/snooze/common/parser"
	"github.com/sirupsen/logrus"
)

type Rule struct {
	If      string   `yaml:"if"`
	GroupBy []string `yaml:"group_by"`
}

type computedRule struct {
	Condition condition.Interface
	Fields    []*field.AlertField
}

var computedRules []*computedRule

func compute(rule *Rule) *computedRule {
	c, err := parser.ParseCondition(rule.If)
	if err != nil {
		log.Fatal(err)
	}
	var fields []*field.AlertField
	for _, f := range rule.GroupBy {
		fi, err := parser.ParseField(f)
		if err != nil {
			log.Fatal(err)
		}
		fields = append(fields, fi)
	}
	return &computedRule{
		Condition: c.Resolve(),
		Fields:    fields,
	}
}

var log *logrus.Entry

func InitRules(rules []*Rule) {
	log = logrus.WithFields(logrus.Fields{"module": "grouping"})
	for _, rule := range rules {
		computedRules = append(computedRules, compute(rule))
	}
}
