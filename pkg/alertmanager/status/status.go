/* Manage the alert status in redis.
Alert status are managed via individual keys, as well as a sorted set
in order to retrieve the expired alerts in O(log(n)).
*/
package status

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	redisv9 "github.com/redis/go-redis/v9"

	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/models"
)

type AlertStatus struct {
	// The ID of the object in Opensearch
	ID string `json:"id"`
	// The key of the alert status in Redis
	Key string `json:"key"`
	LastCheck models.Time `json:"lastCheck"`
	NextCheck models.Time `json:"nextCheck"`
}

// Element of the sorted set. Must be the same upon receiving
// the same alert, so it's a subset of AlertStatus.
type ZAlert struct {
	ID string `json:"id"`
	Key string `json:"key"`
}

const sortedSet = "active-alerts:sorted-set"

func Get(ctx context.Context, key string) (*AlertStatus, bool, error) {
	body, err := redis.Client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return &AlertStatus{}, false, nil
	}

	if err != nil {
		return &AlertStatus{}, false, fmt.Errorf("failed to get redis key '%s': %w", key, err)
	}

	var status *AlertStatus

	if err := json.Unmarshal(body, &status); err != nil {
		return &AlertStatus{}, false, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return status, true, nil
}

func GetExpired(ctx context.Context) ([]ZAlert, error) {
	vals, err := redis.Client.ZRangeByScore(ctx, sortedSet, &redisv9.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
		Count: 1,
		Offset: 0,
	}).Result()
	if err != nil {
		return []ZAlert{}, fmt.Errorf("failed to get the list of expired alert status: %w", err)
	}

	items := []ZAlert{}

	for _, val := range vals {
		var item ZAlert
		if err := json.Unmarshal([]byte(val), &item); err != nil {
			log.Warnf("failed to unmarshal ZAlert `%s` in sorted set: %s", val, err)

			continue
		}

		items = append(items, item)
	}

	return items, nil
}

// Insert or update an alert status.
func Set(ctx context.Context, alertStatus *AlertStatus) error {
	body, err := json.Marshal(alertStatus)
	if err != nil {
		return fmt.Errorf("failed to marshal alert status `%+v`: %w", alertStatus, err)
	}

	zalert := ZAlert{ID: alertStatus.ID, Key: alertStatus.Key}

	setItem, err := json.Marshal(zalert)
	if err != nil {
		return fmt.Errorf("failed to marshal zalert `%+v`: %w", zalert, err)
	}

	pipe := redis.Client.Pipeline()
	pipe.Set(ctx, alertStatus.Key, string(body), 0)
	pipe.ExpireAt(ctx, alertStatus.Key, alertStatus.NextCheck.Time)
	pipe.ZAdd(ctx, sortedSet, redisv9.Z{Member: setItem, Score: float64(time.Now().Unix())})

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to set alert status: %w", err)
	}

	return nil
}

// Remove one or more alert status from the database.
func Delete(ctx context.Context, keys ...interface{}) error {
	if err := redis.Client.ZRem(ctx, sortedSet, keys...).Err(); err != nil {
		return fmt.Errorf("failed to remove alert status from sorted set: %w", err)
	}

	return nil
}
