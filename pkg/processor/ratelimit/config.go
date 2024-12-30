package ratelimit

import (
	"time"

	"github.com/japannext/snooze/pkg/common/lang"
	// "github.com/japannext/snooze/pkg/common/redis"
)

type RateLimit struct {
	// Name of the rate limit. Must be unique (used for key)
	Name string `yaml:"name" json:"name"`

	// An optional condition to only apply the rate limit if the condition match
	If string `yaml:"if" json:"if,omitempty"`

	// The group to group by
	Group string `yaml:"group" json:"group"`

	// The amount of authorized logs during the period
	Burst uint64 `yaml:"burst" json:"burst"`

	// The period used to measure the rate
	Period time.Duration `yaml:"period" json:"period"`

	internal struct {
		condition *lang.Condition
		key       string
	}
}

func (rate *RateLimit) Load() *RateLimit {
	if rate.Name == "" {
		log.Fatalf("ratelimit: missing `name` field!")
	}
	rate.internal.key = rate.Name
	if rate.If != "" {
		condition, err := lang.NewCondition(rate.If)
		if err != nil {
			log.Fatal(err)
		}
		rate.internal.condition = condition
	}

	// limiter := redis_rate.NewLimiter(redis.Client)

	if int64(rate.Period.Seconds()) == 0 {
		log.Fatalf("period should be at least 1 second for rate `%s`", rate.Name)
	}

	return rate
}
