package snooze

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/processor/decision"
	"github.com/japannext/snooze/pkg/processor/metrics"
	redisv9 "github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type Processor struct {
	// storeQ *mq.Pub
}

type Config struct{}

func New(_ Config) (*Processor, error) {
	p := &Processor{}

	return p, nil
}

// Return if the snooze is expired or not.
func snoozeExpired(sz *models.SnoozeLookup) bool {
	now := time.Now()

	return now.After(sz.EndsAt.Time) || now.Before(sz.StartsAt.Time)
}

// Return the snooze entries that match in Redis.
func getSnoozeLookups(ctx context.Context, keys []string) ([]*models.SnoozeLookup, error) {
	lookups := []*models.SnoozeLookup{}
	pipe := redis.Client.Pipeline()
	results := map[string]*redisv9.MapStringStringCmd{}

	log.Debugf("searching for snooze entries: %s", keys)

	for _, key := range keys {
		results[key] = pipe.HGetAll(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if redis.IsError(err) {
		return lookups, fmt.Errorf("failed to execute redis call: %w", err)
	}

	for hash, res := range results {
		entries, err := res.Result()
		if redis.IsNil(err) { // No snooze found, skip.
			continue
		}

		if err != nil {
			log.Warnf("failed get snooze result for hash %s: %s", hash, err)

			continue
		}

		for _, entry := range entries {
			var lookup *models.SnoozeLookup
			if err := json.Unmarshal([]byte(entry), &lookup); err != nil {
				log.Warnf("failed to unmarshal snooze `%s`: %s", entry, err)

				continue
			}

			lookups = append(lookups, lookup)
		}
	}

	return lookups, nil
}

// Check the snooze condition, return true if the log should be skipped
func checkSnoozeCondition(ctx context.Context, sz *models.SnoozeLookup, item *models.Log) bool {
	if sz.If != "" {
		cond, err := lang.NewCondition(sz.If)
		if err != nil {
			log.Errorf("failed to evaluate snooze condition `%s`: %s", sz.If, err)
			return true
		}
		match, err := cond.MatchLog(ctx, item)
		if err != nil {
			log.Errorf("failed to evaluate snooze match `%s`: %s", sz.If, err)
			return true
		}
		if !match {
			return true
		}
	}

	return false
}

func (p *Processor) Process(ctx context.Context, item *models.Log) *decision.Decision {
	ctx, span := otel.Tracer("snooze").Start(ctx, "snooze")
	defer span.End()

	// 1. Get list of possible snooze entries
	keys := []string{}

	for _, group := range item.Groups {
		key := redis.SnoozeKey(group)
		keys = append(keys, key)
	}

	// 2. Check them all in redis
	lookups, err := getSnoozeLookups(ctx, keys)
	if err != nil {
		log.Errorf("failed to fetch snoozes from redis: %s", err)

		return decision.Retry(fmt.Errorf("failed to fetch snoozes from redis: %w", err))
	}

	if len(lookups) == 0 {
		log.Debugf("no snooze entry found")

		return decision.OK()
	}

	// 3. Check if they are still valid.
	for _, lookup := range lookups {
		if skip := checkSnoozeCondition(ctx, lookup, item); skip {
			continue
		}

		if snoozeExpired(lookup) {
			log.Debugf("snooze entry expired. endsAt: %s", lookup.EndsAt)

			continue
		}

		// Change the status to "snoozed"
		ok := item.Status.Change(models.LogSnoozed)
		if ok {
			item.Status.SkipNotification = true
			item.Status.ObjectID = lookup.OpensearchID

			metrics.SnoozedLogs.Inc()
		}

		// Return after first match
		return decision.OK()
	}

	return decision.OK()
}
