package ratelimit

import (
	"fmt"

	"github.com/japannext/snooze/internal/lang"
	"github.com/japannext/snooze/internal/mq"
	"github.com/japannext/snooze/internal/pipeline"
	"github.com/japannext/snooze/internal/utils"
)

type Processor struct {
	ratelimits []Ratelimit
	storeQ     *mq.Pub
}

func New(cfgs []pipeline.Ratelimit) (*Processor, error) {
	p := &Processor{}
	p.storeQ = mq.StorePub()
	duplicates := utils.NewDuplicateChecker()

	for index, cfg := range cfgs {
		if cfg.Name == "" {
			return p, fmt.Errorf("name missing for entry #%d", index+1)
		}

		if err := duplicates.Check(cfg.Name); err != nil {
			return p, fmt.Errorf("duplicate name '%s': %w", cfg.Name, err)
		}

		rl, err := NewRatelimit(cfg)
		if err != nil {
			return p, fmt.Errorf("failed to load '%s': %w", cfg.Name, err)
		}
		p.ratelimits = append(p.ratelimits, rl)
	}

	return p, nil
}

type Ratelimit struct {
	pipeline.Ratelimit
	condition *lang.Condition
	key       string
}

func NewRatelimit(cfg pipeline.Ratelimit) (Ratelimit, error) {
	rl := Ratelimit{}
	if cfg.If != "" {
		condition, err := lang.NewCondition(cfg.If)
		if err != nil {
			return Ratelimit{}, fmt.Errorf("failed to load condition `%s`: %w", cfg.If, err)
		}
		rl.condition = condition
	}

	if int64(cfg.Period.Seconds()) == 0 {
		return Ratelimit{}, fmt.Errorf("period should be at least 1 second for rate `%s`", cfg.Name)
	}

	return rl, nil
}
