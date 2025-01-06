package snooze

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/tracing"
	"github.com/japannext/snooze/pkg/models"
	redisv9 "github.com/redis/go-redis/v9"
)

func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracer.Start(ctx, "snooze")
	defer span.End()

	pipe := redis.Client.Pipeline()
	results := make(map[string]*redisv9.StringCmd)
	hashes := make(map[string]string)
	for _, group := range item.Groups {
		key := fmt.Sprintf("snooze/%s/%s", group.Name, group.Hash)
		results[group.Name] = pipe.Get(ctx, key)
		hashes[group.Name] = group.Hash
	}
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		log.Warnf("failed to get snooze entries: %s", err)
		tracing.Error(span, err)
		return err
	}

	for groupName, res := range results {
		data, err := res.Bytes()
		// No snooze
		if err == redis.Nil {
			continue
		}
		if err != nil {
			log.Warnf("failed get snooze result for group '%s': %s", groupName, err)
			tracing.Error(span, err)
			continue
		}
		var sz *models.Snooze
		if err := json.Unmarshal(data, &sz); err != nil {
			log.Warnf("failed to unmarshal snooze `%s`: %s", data, err)
			tracing.Error(span, err)
			continue
		}
		now := time.Now()
		if now.After(sz.EndsAt.Time) || now.Before(sz.StartsAt.Time) {
			tracing.SetString(span, fmt.Sprintf("snooze.%s:%s", groupName, hashes[groupName]), "ignoring because out of range")
			continue
		}
		if ok := item.Status.Change(models.LogSnoozed); ok {
			item.Status.SkipNotification = true
			item.Status.ObjectID = sz.ID
			snoozedLogs.Inc()
		}
	}

	return nil
}
