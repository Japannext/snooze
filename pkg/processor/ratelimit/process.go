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

// This object is used to represent all variables needed
// to execute the whole process per-rate.
type result struct {
	i0 counter
	i1 counter
	i2 counter

	create *models.Ratelimit
}

// Used to represent a counter used for rate limiting
type counter struct {
	key string
	cmd *redisv9.StringCmd
	value uint64
	hash string
	fields map[string]string
}

// Used to know when a ratelimit is started (to avoid
// starting a ratelimit for every hit)
type lock struct {
	key string
	cmd *redisv9.StatusCmd
}

func (c *counter) extract() error {
	if c.cmd == nil {
		return nil
	}
	value, err := c.cmd.Uint64()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}
	c.value = value
	return nil
}



func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracer.Start(ctx, "ratelimit")
	defer span.End()

	// 1. Filtering ratelimits by condition
	var rr = []*RateLimit{}
	for _, rate := range rates {
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
		rr = append(rr, rate)
	}

	var results = make(map[string]result)

	pipe := redis.Client.Pipeline()
	// 2. Extracting i1/i0 rates
	for _, rate := range rr {
		res := result{}

		period := time.Now().Unix() / int64(rate.Period.Seconds())
		fields, err := lang.ExtractFields(item, rate.internal.fields)
		if err != nil {
			log.Warnf("Error extracting fields: %s", err)
			// skipping[rate.Name] = true
			continue
		}
		hash := utils.ComputeHash(fields)

		res.i0 = counter{
			key: fmt.Sprintf("ratelimits/%s:%s/period:%d", rate.Name, hash, period),
			hash: hash,
			fields: fields,
		}
		res.i1 = counter{
			key: fmt.Sprintf("ratelimits/%s:%s/period:%d", rate.Name, hash, period-1),
			hash: hash,
			fields: fields,
		}
		res.i2 = counter{
			key: fmt.Sprintf("ratelimits/%s:%s/period:%d", rate.Name, hash, period-2),
			hash: hash,
			fields: fields,
		}

		res.i0.cmd = pipe.Get(ctx, res.i0.key)
		res.i1.cmd = pipe.Get(ctx, res.i1.key)
		res.i2.cmd = pipe.Get(ctx, res.i2.key)

		results[rate.Name] = res
	}
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		log.Warnf("failed to get ratelimits: %s", err)
		return err
	}

	// 3. Compute rate info
	pipe = redis.Client.Pipeline()
	for _, rate := range rr {
		res := results[rate.Name]
		i0, i1, i2 := res.i0, res.i1, res.i2

		if err := res.i0.extract(); err != nil {
			log.Warnf("unexpected error while getting i0 value: %s", err)
			continue
		}
		if err := res.i1.extract(); err != nil {
			log.Warnf("unexpected error while getting i1 value: %s", err)
			continue
		}
		if err := res.i2.extract(); err != nil {
			log.Warnf("unexpected error while getting i2 value: %s", err)
			continue
		}

		limit := rate.Burst

		// Rate-limit case
		if i0.value >= limit || i1.value >= limit { // rate limiting
			item.Mute.Drop(fmt.Sprintf("Dropped by ratelimit '%s'", rate.Name))
			rateLimitedLogs.WithLabelValues(rate.Name).Inc()
		}

		// Rate-limit start case
		if i0.value == limit {
			err := storeQ.Publish(ctx, &models.Ratelimit{
				StartsAt: now(),
				Rule: rate.Name,
				Hash: i0.hash,
				Fields: i0.fields,
			})
			if err != nil {
				log.Warnf("failed to publish ratelimit '%s'", rate.Name)
			}
		}

		// Rate-limit stop
		if i2.value > limit && i1.value < limit && i0.value == 0 {
			err := storeQ.WithHeader("action", "update").Publish(ctx, &models.Ratelimit{
				EndsAt: now(),
				Rule: rate.Name,
				Hash: i0.hash,
				Fields: i0.fields,
			})
			if err != nil {
				log.Warnf("failed to update ratelimit '%s'", rate.Name)
			}
		}
	}

	// 4. Increment value for rate
	for _, rate := range rr {
		res := results[rate.Name]
		i0 := res.i0

		pipe.Incr(ctx, i0.key)
		pipe.ExpireNX(ctx, i0.key, rate.Period*3)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		log.Warnf("failed to increment ratelimits: %s", err)
		return err
	}

	return nil
}

func now() uint64 {
	return uint64(time.Now().UnixMilli())
}
