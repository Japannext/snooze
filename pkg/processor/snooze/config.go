package snooze

import (
	"github.com/japannext/snooze/pkg/common/lang"
)

type SnoozeGroup struct {
	Name string `yaml:"name"`
	// An optional condition to only apply the rate limit if the condition match
	If string `yaml:"if"`
	// Fields that will be used to compute a hash for the snooze entry
	Fields []string `yaml:"fields"`

	// A special field representing a map (e.g. identity/labels/source).
	Map string `yaml:"map"`

	internal struct {
		condition *lang.Condition
		fields    []*lang.Field
	}
}
