package ratelimit

import (
	"context"
	"fmt"
	"time"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/redis"
)

func Process(alert *api.Alert) error {
	ctx := context.Background()

	i := time.Now().Unix() / period
	key := fmt.Sprintf("ratelimit/group_hash/%s:%d", alert.GroupHash, i)

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
		alert.Mute.Enabled = true
		alert.Mute.Component = "ratelimit"
		alert.Mute.SkipNotification = true
		alert.Mute.SkipStorage = true
		return fmt.Errorf("ratelimit: immediate stop")
	}

	return nil
}
