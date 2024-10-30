package grouping

import (
	"slices"

	"github.com/japannext/snooze/pkg/common/lang"
)

// Group logs by fields. Groups can then be used
// for rate-limiting, snooze, and UI search.
type Grouping struct {
	Name	string   `yaml:"name"`
	If      string   `yaml:"if"`
	// Mutually exclusive with `group_by_map`.
	GroupBy []string `yaml:"group_by"`
	// Mutually exclusive with `group_by`.
	GroupByMap string `yaml:"group_by_map"`

	// A string to help formatting the group.
	FormatLabels string `yaml:"format_labels"`

	internal struct {
		condition *lang.Condition
		fields []*lang.Field
	}
}

var GROUP_BY_MAP = []string{"source", "identity", "labels"}

func (group *Grouping) Load() *Grouping {
	if group.If != "" {
		condition, err := lang.NewCondition(group.If)
		if err != nil {
			log.Fatal(err)
		}
		group.internal.condition = condition
	}
	if len(group.GroupBy) > 0 {
		fields, err := lang.NewFields(group.GroupBy)
		if err != nil {
			log.Fatal(err)
		}
		group.internal.fields = fields
	}

	if group.GroupByMap != "" && len(group.GroupBy) != 0 {
		log.Fatalf("group_by and group_by_map are mutually exclusive")
	}

	if group.GroupByMap != "" && !slices.Contains(GROUP_BY_MAP, group.GroupByMap) {
		log.Fatalf("group_by_map='%s' is invalid. Allowed values: source, identity, labels", group.GroupByMap)
	}

	return group
}

/*
groupings:
- name: by-source
  map: source
- name: by-host
  if: 'has(identity["host"])'
  fields: ['identity.host']
*/
