package snooze

import (
	"context"
	"fmt"

	redisv9 "github.com/redis/go-redis/v9"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/tracing"
)

func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracer.Start(ctx, "snooze")
	defer span.End()

	pipe := redis.Client.Pipeline()
	results := make(map[string]*redisv9.StringCmd)
	for _, group := range item.Groups {
		key := fmt.Sprintf("snooze/%s/%s", group.Name, group.Hash)
		results[group.Name] = pipe.Get(ctx, key)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		log.Warnf("failed to get snooze entries: %s", err)
		tracing.Error(span, err)
		return err
	}

	for groupName, res := range results {
		reason, err := res.Result()
		// No snooze
		if err == redis.Nil {
			continue
		}
		if err != nil {
			log.Warnf("failed get snooze result for group '%s': %s", groupName, err)
			tracing.Error(span, err)
			continue
		}
		item.Mute.Silence(fmt.Sprintf("Snoozed by group '%s': %s", groupName, reason))
	}

	return nil
}
