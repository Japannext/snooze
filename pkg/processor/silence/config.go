package silence

import (
	"fmt"

	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/common/schedule"
)

type Silence struct {
	Name        string             `json:"name"                  yaml:"name"`
	Description string             `json:"description,omitempty" yaml:"description"`
	If          string             `json:"if"                    yaml:"if"`
	Schedule    *schedule.Schedule `yaml:",inline"`
	Drop        bool               `json:"drop"                  yaml:"drop"`

	internal struct {
		condition *lang.Condition
	}
}

func (s *Silence) String() string {
	if s.Name != "" {
		return s.Name
	}
	return fmt.Sprintf("`%s`, ", s.If)
}

func (s *Silence) Load() {
	var err error
	s.internal.condition, err = lang.NewCondition(s.If)
	if err != nil {
		log.Fatal(err)
	}

	if s.Schedule == nil {
		s.Schedule = schedule.Default()
	}
	s.Schedule.Load()
}
