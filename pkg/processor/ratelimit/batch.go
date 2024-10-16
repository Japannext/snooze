package ratelimit

import (
	"context"
	"fmt"
	"time"

	redisv9 "github.com/redis/go-redis/v9"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/common/utils"
)

func newMap[T any](n int) []map[string]T {
	mm := make([]map[string]T, n)
	for i, _ := range mm {
		mm[i] = make(map[string]T)
	}
	return mm
}

// Process a batch of items
func Batch(ctx context.Context, items []*models.Log) error {
	ctx, span := tracer.Start(ctx, "ratelimit")
	defer span.End()

	n := len(items)

	// This stores the variables per-rate/per-item,
	// but because we batch, we need to save these
	// variables between loops
	previousKeys := make(map[string]string)
	currentKeys := make(map[string]string)
	previousCmds := map[string]*redisv9.StringCmd{}
	currentCmds := map[string]*redisv9.StringCmd{}

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
			k := fmt.Sprintf("%s:%d", rate.Name, i)
			previousKeys[k] = fmt.Sprintf("ratelimits/%s:%s/period:%d", rate.Name, hash, period)
			currentKeys[k] = fmt.Sprintf("ratelimits/%s:%s/period:%d", rate.Name, hash, period-1)
			previousCmds[k] = pipe.Get(ctx, previousKeys[k])
			currentCmds[k] = pipe.Get(ctx, currentKeys[k])
		}
	}
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		log.Warnf("failed to get ratelimits for %d logs: %s", n, err)
		return err
	}
	for _, rate := range rates {
		// Verify burst is not passed
		for i, item := range items {
			k := fmt.Sprintf("%s:%d", rate.Name, i)
			previous, err := extractInteger(previousCmds, k)
			if err != nil {
				log.Warnf("unexpected error while getting previous value: %s", err)
				continue
			}
			current, err := extractInteger(currentCmds, k)
			if err != nil {
				log.Warnf("unexpected error while getting current value: %s", err)
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
			k := fmt.Sprintf("%s:%d", rate.Name, i)
			pipe.Incr(ctx, currentKeys[k])
			pipe.ExpireNX(ctx, currentKeys[k], rate.Period*3)
		}
	}
	if _, err := pipe.Exec(ctx); err != nil {
		log.Warnf("failed to increment ratelimits for %d logs!", n)
		return err
	}

	return nil
}

// Extract an integer, and handle all edge cases
func extractInteger(cmds map[string]*redisv9.StringCmd, key string) (uint64, error) {
	cmd, found := cmds[key]
	if !found {
		return 0, nil
	}
	value, err := cmd.Uint64()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return value, nil

}
