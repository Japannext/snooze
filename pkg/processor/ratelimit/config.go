package ratelimit

import (
	"time"

	"github.com/japannext/snooze/pkg/common/lang"
)

type RateLimit struct {
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
