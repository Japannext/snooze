package ratelimit

import (
  "context"
  "fmt"
  "time"
  "strconv"

  log "github.com/sirupsen/logrus"
  "github.com/japannext/snooze/common/api/v2"
  "github.com/japannext/snooze/common/redis"
)

type RateLimiter struct {
  Burst int
  Period time.Duration
}

func (r *RateLimiter) Process(ctx context.Context, item v2.Alert) ([]v2.Alert, error) {

  var out []v2.Alert
  var keys []string
  // Unique integer value per chunk of Period.
  i := time.Now().Unix() / int64(r.Period.Seconds())
  key := fmt.Sprintf("ratelimit/group_hash/%s:%d", item.GroupHash, i)

  rdb := redis.Client
  vals, err := rdb.Get(ctx, keys).Result()
  if err == redis.Nil {
    // Not found
    rdb.Set(ctx, key, 1)
  }
  if err != nil {
    return out, err
  }

  pipe := rdb.Pipeline()
  for i, events := range events {
    key := keys[i]
    val := vals[i]
    v, ok := val.(int64)
    if !ok {
      // Weird, but let's remediate.
      log.Warnf("[ratelimiter] Non-integer value for key '%s': %v+", key, val)
      pipe.Set(ctx, key, 0)
    } else {
      pipe.Incr(ctx, key)
      pipe.ExpireNX(ctx, key, r.Period)
    }
  }
  cmds, err := pipe.Exec(ctx)
  if err != nil {
    return out, err
  }

  for i, val := range vals {
    key = keys[i]
    v, ok := val.(int64)
    if !ok {
      // ???
      continue
    }
    if v <= r.Burst {
      out = append(out, )
    }
  }

  for _, event := range events {
    key := fmt.Sprintf("ratelimit/group_hash/%s:%d", event.GroupHash, i)
    c, err := rdb.Get(ctx, key).Uint64()
    if err != nil {
      continue
    }
    if c > r.Burst {
      // rate limit ?
      // notification ?
    }
    err := rdb.TxPipelinex(ctx, func(pipe.Pipeliner) error {
      if err := pipe.Incr(ctx, key); err != nil {
        return err
      }
      if err := pipe.ExpireNX(ctx, key, r.Period); err != nil {
        return err
      }
      return nil
    })
  }
}

type RateLimitView struct {
  Firing bool
  Count int64
  Invalid bool
}

// Example for the web interface ?
func (r *RateLimiter) IsRateLimited(ctx context.Context, events []v2.Alert) ([]RateLimitView, error) {
  var keys []string
  i := time.Now().Unix() / r.Period
  for _, event := range events {
    keys = append(keys, fmt.Sprintf("ratelimit/group_hash/%s:%d", event.GroupHash, i))
  }

  rdb := redis.DB
  vals, err := rdb.MGet(ctx, keys).Result()
  if err != nil {
    return []RateLimitView{}, err
  }

  var views []RateLimitView
  for _, val := range vals {
    if v, ok := val.(int64); ok {
      view := RateLimitView{Count: v}
      if v > r.Burst {
        view.Firing = true
      }
    } else {
      view := RateLimitView{Invalid: true}
    }
    views = append(views, view)
  }
  return views, nil
}
