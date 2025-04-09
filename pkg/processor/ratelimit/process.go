package ratelimit

import (
	"context"
	"fmt"

	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/utils"
	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/processor/decision"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

// Number of retries for the redis CL.THROTTLE call
const maxRetries = 1

// Return only the ratelimits that match the log (if conditions are set)
func (p *Processor) filteredRatelimits(ctx context.Context, item *models.Log) []Ratelimit {
	ratelimits := []Ratelimit{}
	for _, rate := range p.ratelimits {
		if rate.internal.condition != nil {
			match, err := rate.internal.condition.MatchLog(ctx, item)
			if err != nil {
				log.Warnf("[rate=%s]failed to match condition `%s`: %s", rate.Name, rate.internal.condition, err)
				continue
			}
			if !match {
				continue
			}
		}
		if _, ok := utils.GetGroup(item, rate.Group); !ok {
			continue
		}
		ratelimits = append(ratelimits, rate)
	}

	return ratelimits
}

func getRatelimitKey(gr *models.Group) string {
	return fmt.Sprintf("ratelimit:%s:%s", gr.Name, gr.Hash)
}

func (p *Processor) Process(ctx context.Context, item *models.Log) *decision.Decision {
	ctx, span := otel.Tracer("snooze").Start(ctx, "ratelimit")
	defer span.End()

	// 1. Filtering ratelimits by condition
	ratelimits := p.filteredRatelimits(ctx, item)

	for _, rate := range ratelimits {
		gr, _ := utils.GetGroup(item, rate.Group)
		key := getRatelimitKey(gr)

		throttle, err := redis.Client.CLThrottle(ctx, key, maxRetries, rate.CountPerPeriod, rate.Period, 1)
		if err != nil {
			log.Warnf("failed to execute ratelimit %s: %s", key, err)
			continue
		}

		if !throttle.Allowed {
			if err := UpsertStatus(ctx, rate, gr, throttle); err != nil {
				log.Warnf("failed to upsert ratelimit: %s", err)
				return decision.OK()
			}
			return decision.Abort(fmt.Errorf("ratelimited by `%s`: retry in %s", rate.Name, throttle.RetryAfter))
		}
	}

	return decision.OK()
}
