package redis

import (
	"context"
	"fmt"
	"time"
)

type ThrottleState struct {
	Allowed    bool
	Limit      int64
	Remaining  int64
	RetryAfter time.Duration
	Reset      time.Duration
}

// Wrapper of redis-cell plugin / dragonflydb command
func (client *RedisClient) CLThrottle(ctx context.Context, key string, maxBurst, countPerPeriod int64, period time.Duration, quantity int64) (*ThrottleState, error) {
	res, err := client.Do(ctx, "CL.THROTTLE", key, maxBurst, countPerPeriod, int64(period.Seconds()), quantity).Int64Slice()
	if err != nil {
		return &ThrottleState{}, fmt.Errorf("error fetching CL.THROTTLE %s: %w", key, err)
	}

	if len(res) < 5 {
		return &ThrottleState{}, fmt.Errorf("invalid returned value for CL.THROTTLE (less than 5 integers returned)")
	}

	return &ThrottleState{
		Allowed:    (res[0] == 0),
		Limit:      res[1],
		Remaining:  res[2],
		RetryAfter: time.Duration(res[3]) * time.Second,
		Reset:      time.Duration(res[4]) * time.Second,
	}, nil
}
