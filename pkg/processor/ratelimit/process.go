package ratelimit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/japannext/snooze/pkg/common/opensearch/format"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/tracing"
	"github.com/japannext/snooze/pkg/common/utils"
	"github.com/japannext/snooze/pkg/models"
	redisv9 "github.com/redis/go-redis/v9"
)

// This object is used to represent all variables needed
// to execute the whole process per-rate.
type result struct {
	i0 counter
	i1 counter
	i2 counter

	lock lock
}

// Used to represent a counter used for rate limiting.
type counter struct {
	key    string
	cmd    *redisv9.StringCmd
	value  uint64
	hash   string
	fields map[string]string
}

// Used to know when a ratelimit is started (to avoid
// starting a ratelimit for every hit).
type lock struct {
	key string
	cmd *redisv9.StatusCmd
}

func (c *counter) extract() error {
	if c.cmd == nil {
		return fmt.Errorf("No cmd provided (nil) for '%s'", c.key)
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
	rr := []*RateLimit{}
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
		if _, ok := utils.GetGroup(item, rate.Group); !ok {
			continue
		}
		rr = append(rr, rate)
	}

	results := make(map[string]result)

	pipe := redis.Client.Pipeline()
	// 2. Extracting i1/i0 rates
	for _, rate := range rr {
		res := result{}

		period := time.Now().Unix() / int64(rate.Period.Seconds())
		group, _ := utils.GetGroup(item, rate.Group)
		hash := group.Hash
		fields := group.Labels

		res.i0 = counter{
			key:    fmt.Sprintf("ratelimits/%s:%s/period:%d", rate.Name, hash, period),
			hash:   hash,
			fields: fields,
		}
		res.i1 = counter{
			key:    fmt.Sprintf("ratelimits/%s:%s/period:%d", rate.Name, hash, period-1),
			hash:   hash,
			fields: fields,
		}
		res.i2 = counter{
			key:    fmt.Sprintf("ratelimits/%s:%s/period:%d", rate.Name, hash, period-2),
			hash:   hash,
			fields: fields,
		}
		res.lock = lock{
			key: fmt.Sprintf("ratelimits/%s:%s/lock", rate.Name, hash),
		}

		res.i0.cmd = pipe.Get(ctx, res.i0.key)
		res.i1.cmd = pipe.Get(ctx, res.i1.key)
		res.i2.cmd = pipe.Get(ctx, res.i2.key)
		res.lock.cmd = pipe.SetArgs(ctx, res.lock.key, uuid.NewString(), redisv9.SetArgs{Get: true, Mode: "NX"})

		results[rate.Name] = res
	}
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		log.Warnf("failed to get ratelimits: %s", err)
		tracing.Error(span, err)
		return err
	}

	// 3. Compute rate info
	for _, rate := range rr {
		res := results[rate.Name]
		i0, i1, i2, lock := res.i0, res.i1, res.i2, res.lock

		if err := res.i0.extract(); err != nil {
			log.Warnf("unexpected error while getting i0 value: %s", err)
			tracing.Error(span, err)
			continue
		}
		if err := res.i1.extract(); err != nil {
			log.Warnf("unexpected error while getting i1 value: %s", err)
			tracing.Error(span, err)
			continue
		}
		if err := res.i2.extract(); err != nil {
			log.Warnf("unexpected error while getting i2 value: %s", err)
			tracing.Error(span, err)
			continue
		}

		tracing.SetInt(span, fmt.Sprintf("ratelimit.%s.i2", rate.Name), int(i2.value))
		tracing.SetInt(span, fmt.Sprintf("ratelimit.%s.i1", rate.Name), int(i1.value))
		tracing.SetInt(span, fmt.Sprintf("ratelimit.%s.i0", rate.Name), int(i0.value))

		limit := rate.Burst

		// Rate-limit case
		if i0.value >= limit || i1.value >= limit { // rate limiting
			if ok := item.Status.Change(models.LogRatelimited); ok {
				item.Status.Reason = fmt.Sprintf("ratelimited by '%s'", rate.Name)
				item.Status.SkipNotification = true
				item.Status.SkipStorage = true
				rateLimitedLogs.WithLabelValues(rate.Name).Inc()
			}
			tracing.SetString(span, fmt.Sprintf("ratelimit.%s.decision", rate.Name), "drop")
		} else {
			tracing.SetString(span, fmt.Sprintf("ratelimit.%s.decision", rate.Name), "store")
		}

		id := lock.cmd.String()

		// Rate-limit start case
		if i0.value >= limit || id == "" {
			item := &models.Ratelimit{
				StartsAt: now(),
				Rule:     rate.Name,
				Hash:     i0.hash,
				Fields:   i0.fields,
			}
			err := storeQ.PublishData(ctx, &format.Create{
				Index: models.RatelimitIndex,
				Item:  item,
			})
			if err != nil {
				log.Warnf("failed to publish ratelimit '%s'", rate.Name)
				tracing.Error(span, err)
			}
		}

		// Rate-limit stop
		if i2.value > limit && i1.value < limit && i0.value == 0 && lock.cmd.String() != "" {
			id := lock.cmd.String()

			rt := &models.Ratelimit{EndsAt: now()}
			data, err := json.Marshal(rt)
			if err != nil {
				log.Warnf("failed to marshal ratelimit %+v: %s", rt, err)
				tracing.Error(span, err)
				continue
			}
			err = storeQ.PublishData(ctx, &format.Update{
				Index: models.RatelimitIndex,
				ID:    id,
				Doc:   data,
			})
			if err != nil {
				log.Warnf("failed to update ratelimit '%s'", rate.Name)
				tracing.Error(span, err)
			}
		}
	}

	// 4. Increment value for rate
	var incrCmds []*redisv9.IntCmd
	var expireCmds []*redisv9.BoolCmd
	pipe = redis.Client.Pipeline()
	for _, rate := range rr {
		i0 := results[rate.Name].i0

		incrCmd := pipe.Incr(ctx, i0.key)
		expireCmd := pipe.ExpireNX(ctx, i0.key, rate.Period*3)
		incrCmds = append(incrCmds, incrCmd)
		expireCmds = append(expireCmds, expireCmd)
		tracing.SetString(span, fmt.Sprintf("ratelimit.%s.incr_key", rate.Name), i0.key)
	}
	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		log.Warnf("failed to increment ratelimits: %s", err)
		tracing.Error(span, err)
		return err
	}
	for _, incrCmd := range incrCmds {
		if err := incrCmd.Err(); err != nil {
			tracing.Error(span, err)
			log.Warnf("failed to increment in redis: %s", err)
		}
	}
	for _, expireCmd := range expireCmds {
		if err := expireCmd.Err(); err != nil {
			tracing.Error(span, err)
			log.Warnf("failed to set expire in redis: %s", err)
		}
	}

	return nil
}

func now() uint64 {
	return uint64(time.Now().UnixMilli())
}
