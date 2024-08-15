package ratelimit

import (
	"context"
	"fmt"
	"time"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/redis"
)

func Process(item *api.Log) error {
	ctx := context.Background()

	i := time.Now().Unix() / period
	key := fmt.Sprintf("ratelimit/group_hash/%s:%d", item.GroupHash, i)

	// Test for existence and set if needed
	v, err := rdb.Get(ctx, key).Int64()
	if err == redis.Nil {
		return rdb.SetNX(ctx, key, 1, rateLimit.Period).Err()
	}
	if err != nil {
		return err
	}

	// Increment and update expire if needed
	pipe := rdb.Pipeline()
	pipe.Incr(ctx, key)
	pipe.ExpireNX(ctx, key, rateLimit.Period)

	if v > burst {
		item.Mute.Drop("Dropped by ratelimit")
		return nil
	}

	return nil
}
