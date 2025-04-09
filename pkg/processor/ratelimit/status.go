package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/models"
)

func statusKey(ruleName, groupName, hash string) string {
	return fmt.Sprintf("ratelimit:status:%s:%s", ruleName, groupName, hash)
}

// Create a new rate-limit status in Redis
func newStatus(ctx context.Context, key string, rule Ratelimit, gr *models.Group, throttle *redis.ThrottleState) error {
	pipe := redis.Client.Pipeline()

	now := uint64(time.Now().UnixMilli())
	args := []interface{}{
		"active", true,
		"startsAt", now,
		"lastHit", now,
		"rule", rule.Name,
		"hash", gr.Hash,
		"hits", 1,
	}
	pipe.HSet(ctx, key, args...)
	pipe.Expire(ctx, key, throttle.RetryAfter)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to set ratelimit status: %w", err)
	}

	return nil
}

// Add the lastHit and hit number to the ratelimit status.
func updateStatus(ctx context.Context, key string, throttle *redis.ThrottleState) error {
	pipe := redis.Client.Pipeline()
	pipe.HSet(ctx, key, "lastHit", uint64(time.Now().UnixMilli()))
	pipe.HIncrBy(ctx, key, "hits", 1)
	pipe.Expire(ctx, key, throttle.RetryAfter)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	return nil
}

// Create or update the status
func UpsertStatus(ctx context.Context, rule Ratelimit, gr *models.Group, throttle *redis.ThrottleState) error {
	key := statusKey(rule.Name, gr.Name, gr.Hash)
	res, err := redis.Client.Exists(ctx, key).Result()
	if err != nil {

	}

	// Key doesn't exists
	if res == 0 {
		err = newStatus(ctx, key, rule, gr, throttle)
	} else {
		err = updateStatus(ctx, key, throttle)
	}
	if err != nil {
		return fmt.Errorf("failed to upsert: %w", err)
	}

	return nil
}

// Close the status when it has been decided that this ratelimit was over
func CloseStatus(ctx context.Context, rule Ratelimit, gr *models.Group) error {
	key := statusKey(rule.Name, gr.Name, gr.Hash)
	res, err := redis.Client.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to check if %s exists: %w", key, err)
	}
	if res == 0 {
		return nil
	}

	if err := redis.Client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}

	return nil
}
