package ratelimit

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/common/utils"
)

func Process(item *api.Log) error {
	ctx := context.Background()

	for _, rate := range rates {
		log.Debugf("Evaluating ratelimit '%s'", rate.Name)
		// Condition
		if rate.internal.condition != nil {
			match, err := rate.internal.condition.MatchLog(ctx, item)
			if err != nil {
				return err
			}
			if !match {
				continue
			}
		}

		// Fields hash
		fields, err := lang.ExtractFields(item, rate.internal.fields)
		if err != nil {
			return err
		}
		hash := hex.EncodeToString(utils.ComputeHash(fields))

		// Check redis
		period := time.Now().Unix() / int64(rate.Period.Seconds())

		previousKey := fmt.Sprintf("ratelimits/%s:%s/period:%d", rate.Name, hash, period)
		previous, err := redis.Client.Get(ctx, previousKey).Uint64()
		if err != nil && err != redis.Nil {
			log.Warnf("error fetching previous value: %s", err)
		}

		currentKey := fmt.Sprintf("ratelimits/%s:%s/period:%d", rate.Name, hash, period-1)
		current, err := redis.Client.Get(ctx, currentKey).Uint64()
		if err != nil && err != redis.Nil {
			log.Warnf("error fetching current value: %s", err)
		}

		if previous > rate.Burst || current > rate.Burst {
			log.Warnf("dropped by ratelimit `%s`", rate.Name)
			item.Mute.Drop(fmt.Sprintf("Dropped by ratelimit '%s'", rate.Name))
			rateLimitedLogs.WithLabelValues(rate.Name).Inc()
		}

		// Increment and update expire if needed
		pipe := redis.Client.Pipeline()
		pipe.Incr(ctx, currentKey)
		pipe.ExpireNX(ctx, currentKey, rate.Period*2)
		if _, err := pipe.Exec(ctx); err != nil {
			log.Warnf("failed to increment ratelimit %s", rate.Name)
			continue
		}
		log.Debugf("Incremented %s for %s", currentKey, rate.Period*2)
	}

	return nil
}
