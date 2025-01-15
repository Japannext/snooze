package ratelimit

import (
	"fmt"
	"time"

	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/utils"
)

type Processor struct {
	ratelimits []Ratelimit
	storeQ     *mq.Pub
}

type Config struct {
	Ratelimits []Ratelimit `json:"ratelimits" yaml:"ratelimits"`
}

func New(cfg Config) (*Processor, error) {
	p := &Processor{}
	p.storeQ = mq.StorePub()
	duplicates := utils.NewDuplicateChecker()

	for index, rate := range cfg.Ratelimits {
		if rate.Name == "" {
			return p, fmt.Errorf("name missing for entry #%d", index+1)
		}

		if err := duplicates.Check(rate.Name); err != nil {
			return p, fmt.Errorf("duplicate name '%s': %w", rate.Name, err)
		}

		if err := rate.Load(); err != nil {
			return p, fmt.Errorf("failed to load '%s': %w", rate.Name, err)
		}
	}

	return p, nil
}

type Ratelimit struct {
	// Name of the rate limit. Must be unique (used for key)
	Name string `json:"name" yaml:"name"`

	// An optional condition to only apply the rate limit if the condition match
	If string `json:"if,omitempty" yaml:"if"`

	// The group to group by
	Group string `json:"group" yaml:"group"`

	// The amount of authorized logs during the period
	Burst uint64 `json:"burst" yaml:"burst"`

	// The period used to measure the rate
	Period time.Duration `json:"period" yaml:"period"`

	internal struct {
		condition *lang.Condition
		key       string
	}
}

func (rate *Ratelimit) Load() error {
	rate.internal.key = rate.Name
	if rate.If != "" {
		condition, err := lang.NewCondition(rate.If)
		if err != nil {
			return fmt.Errorf("failed to load condition `%s`: %w", rate.If, err)
		}
		rate.internal.condition = condition
	}

	if int64(rate.Period.Seconds()) == 0 {
		return fmt.Errorf("period should be at least 1 second for rate `%s`", rate.Name)
	}

	return nil
}
