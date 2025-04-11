package snooze

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	// "github.com/japannext/snooze/pkg/common/mq"
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
func snoozeExpired(sz *models.Snooze) bool {
	now := time.Now()

	return now.After(sz.EndsAt.Time) || now.Before(sz.StartsAt.Time)
}

// Return the snooze entries that match in Redis.
func getSnoozes(ctx context.Context, keys []string) ([]*models.Snooze, error) {
	snoozes := []*models.Snooze{}
	pipe := redis.Client.Pipeline()
	results := map[string]*redisv9.StringCmd{}

	log.Debugf("searching for snooze entries: %s", keys)

	for _, key := range keys {
		results[key] = pipe.Get(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if redis.IsError(err) {
		return snoozes, fmt.Errorf("failed to execute redis call: %w", err)
	}

	for hash, res := range results {
		data, err := res.Bytes()
		if redis.IsNil(err) { // No snooze found, skip.
			continue
		}

		if err != nil {
			log.Warnf("failed get snooze result for hash %s: %s", hash, err)

			continue
		}

		var sz *models.Snooze

		if err := json.Unmarshal(data, &sz); err != nil {
			log.Warnf("failed to unmarshal snooze `%s`: %s", data, err)

			continue
		}

		snoozes = append(snoozes, sz)
	}

	return snoozes, nil
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
	snoozes, err := getSnoozes(ctx, keys)
	if err != nil {
		log.Errorf("failed to fetch snoozes from redis: %s", err)

		return decision.Retry(fmt.Errorf("failed to fetch snoozes from redis: %w", err))
	}

	if len(snoozes) == 0 {
		log.Debugf("no snooze entry found")

		return decision.OK()
	}

	// 3. Check if they are still valid.
	for _, sz := range snoozes {
		if snoozeExpired(sz) {
			log.Debugf("snooze entry expired. endsAt: %s", sz.EndsAt)

			continue
		}

		// Change the status to "snoozed"
		ok := item.Status.Change(models.LogSnoozed)
		if ok {
			item.Status.SkipNotification = true
			item.Status.ObjectID = sz.ID

			metrics.SnoozedLogs.Inc()
		}

		// Return after first match
		return decision.OK()
	}

	return decision.OK()
}
