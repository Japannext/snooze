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

func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracer.Start(ctx, "ratelimit")
	defer span.End()

	var (
		previousKeys = make(map[string]string)
		currentKeys = make(map[string]string)
		previousCmds = make(map[string]*redisv9.StringCmd)
		currentCmds = make(map[string]*redisv9.StringCmd)
		skipping = make(map[string]bool)
	)

	// 1. Skip if condition doesn't match
	for _, rate := range rates {
		if rate.internal.condition != nil {
			match, err := rate.internal.condition.MatchLog(ctx, item)
			if err != nil {
				log.Warnf("[rate=%s]failed to match condition `%s`: %s", rate.Name, rate.internal.condition, err)
				skipping[rate.Name] = true
				continue
			}
			if !match {
				skipping[rate.Name] = true
			}
		}
	}

	pipe := redis.Client.Pipeline()
	// 2. Extracting previous/current rates
	for _, rate := range rates {
		if _, skip := skipping[rate.Name]; skip {
			continue
		}
		period := time.Now().Unix() / int64(rate.Period.Seconds())
		log.Debugf("time = %d, config.Period.Seconds() = %d, period = %d", time.Now().Unix(), int64(rate.Period.Seconds()), period)
		fields, err := lang.ExtractFields(item, rate.internal.fields)
		if err != nil {
			log.Warnf("Error extracting fields: %s", err)
			skipping[rate.Name] = true
			continue
		}
		hash := utils.ComputeHash(fields)
		previousKeys[rate.Name] = fmt.Sprintf("ratelimits/%s:%s/period:%d", rate.Name, hash, period)
		currentKeys[rate.Name] = fmt.Sprintf("ratelimits/%s:%s/period:%d", rate.Name, hash, period-1)
		previousCmds[rate.Name] = pipe.Get(ctx, previousKeys[rate.Name])
		currentCmds[rate.Name] = pipe.Get(ctx, currentKeys[rate.Name])
	}
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		log.Warnf("failed to get ratelimits: %s", err)
		return err
	}

	// 3. Compute rate info
	for _, rate := range rates {
		if _, skip := skipping[rate.Name]; skip {
			continue
		}
		previous, err := extractInteger(previousCmds, rate.Name)
		if err != nil {
			log.Warnf("unexpected error while getting previous value: %s", err)
			continue
		}
		current, err := extractInteger(currentCmds, rate.Name)
		if err != nil {
			log.Warnf("unexpected error while getting current value: %s", err)
			continue
		}
		// Checking rate
		if previous > rate.Burst || current > rate.Burst {
			item.Mute.Drop(fmt.Sprintf("Dropped by ratelimit '%s'", rate.Name))
			rateLimitedLogs.WithLabelValues(rate.Name).Inc()
		}
	}

	// 4. Increment value for rate
	pipe = redis.Client.Pipeline()
	for _, rate := range rates {
		if _, skip := skipping[rate.Name]; skip {
			continue
		}
		pipe.Incr(ctx, currentKeys[rate.Name])
		pipe.ExpireNX(ctx, currentKeys[rate.Name], rate.Period*3)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		log.Warnf("failed to increment ratelimits: %s", err)
		return err
	}

	return nil
}

/*
func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracer.Start(ctx, "ratelimit")
	defer span.End()
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
		hash := utils.ComputeHash(fields)

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
*/
