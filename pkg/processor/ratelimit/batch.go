package ratelimit

import (
	"context"
	"fmt"
	"time"

	redisv9 "github.com/redis/go-redis/v9"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/common/utils"
)

// Process a batch of items
func Batch(items []*api.Log) {
	ctx := context.Background()

	n := len(items)

	// This stores the variables per-rate/per-item,
	// but because we batch, we need to save these
	// variables between loops
	previousKeys := make([]map[string]string, n)
	currentKeys := make([]map[string]string, n)
	previousCmds := make([]map[string]*redisv9.StringCmd, n)
	currentCmds := make([]map[string]*redisv9.StringCmd, n)

	// First redis pipe: get all previous/current values
	pipe := redis.Client.Pipeline()
	for _, rate := range rates {
		period := time.Now().Unix() / int64(rate.Period.Seconds())
		for i, item := range items {
			// Fields hash
			fields, err := lang.ExtractFields(item, rate.internal.fields)
			if err != nil {
				log.Warnf("Error extracting fields: %s", err)
				continue
			}
			hash := utils.ComputeHash(fields)
			previousKeys[i][rate.Name] = fmt.Sprintf("ratelimits/%s:%s/period:%d", rate.Name, hash, period)
			currentKeys[i][rate.Name] = fmt.Sprintf("ratelimits/%s:%s/period:%d", rate.Name, hash, period-1)
			previousCmds[i][rate.Name] = pipe.Get(ctx, previousKeys[i][rate.Name])
			currentCmds[i][rate.Name] = pipe.Get(ctx, currentKeys[i][rate.Name])
		}
	}
	if _, err := pipe.Exec(ctx); err != nil {
		log.Warnf("failed to get ratelimits for %d logs!", n)
		return
	}
	for _, rate := range rates {
		// Verify burst is not passed
		for i, item := range items {
			previous, err := previousCmds[i][rate.Name].Uint64()
			if err != nil {
				log.Warnf("invalid format for previous ratelimit")
				continue
			}
			current, err := currentCmds[i][rate.Name].Uint64()
			if err != nil {
				log.Warnf("invalid format for current ratelimit")
				continue
			}
			if previous > rate.Burst || current > rate.Burst {
				log.Warnf("dropped by ratelimit `%s`", rate.Name)
				item.Mute.Drop(fmt.Sprintf("Dropped by ratelimit '%s'", rate.Name))
				rateLimitedLogs.WithLabelValues(rate.Name).Inc()
			}
		}
	}

	// Second pipe: increment current value
	pipe = redis.Client.Pipeline()
	for _, rate := range rates {
		for i, _ := range items {
			pipe.Incr(ctx, currentKeys[i][rate.Name])
			pipe.ExpireNX(ctx, currentKeys[i][rate.Name], rate.Period*2)
		}
	}
	if _, err := pipe.Exec(ctx); err != nil {
		log.Warnf("failed to increment ratelimits for %d logs!", n)
		return
	}
}
