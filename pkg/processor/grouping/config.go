package grouping

import (
	"github.com/japannext/snooze/pkg/common/lang"
)

type Grouping struct {
	Name	string   `yaml:"name"`
	If      string   `yaml:"if"`
	GroupBy []string `yaml:"group_by"`

	internal struct {
		condition *lang.Condition
		fields []*lang.Field
	}
}

func (group *Grouping) Load() *Grouping {
	if group.If != "" {
		condition, err := lang.NewCondition(group.If)
		if err != nil {
			log.Fatal(err)
		}
		group.internal.condition = condition
	}
	fields, err := lang.NewFields(group.GroupBy)
	if err != nil {
		log.Fatal(err)
	}
	group.internal.fields = fields

	return group
}
