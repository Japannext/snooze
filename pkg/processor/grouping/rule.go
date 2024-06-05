package grouping

import (
	"github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/common/lang"
)

type Rule struct {
	If      string   `yaml:"if"`
	GroupBy []string `yaml:"group_by"`
}

type computedRule struct {
	Condition *lang.Condition
	Fields    []*lang.Field
}

var computedRules []*computedRule

func compute(rule *Rule) *computedRule {
	condition, err := lang.NewCondition(rule.If)
	if err != nil {
		log.Fatal(err)
	}
	var fields []*lang.Field
	for _, groupby := range rule.GroupBy {
		f, err := lang.NewField(groupby)
		if err != nil {
			log.Fatal(err)
		}
		fields = append(fields, f)
	}
	return &computedRule{
		Condition: condition,
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
