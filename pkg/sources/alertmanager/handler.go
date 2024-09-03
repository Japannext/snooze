package alertmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/utils"
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

const ALERTMANAGER_KIND = "alertmanager"
const TIME_LAYOUT = "2006-01-02T03:04:05.999999999Z"

func getAlertStatus(ctx context.Context, key string) (*api.AlertStatus, error) {
	body, err := redis.Client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if err == redis.Nil {
		return nil, nil
	}
	var status *api.AlertStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, err
	}
	return status, nil
}

func setAlertStatus(ctx context.Context, key string, status *api.AlertStatus) error {
	body, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal alert status `%+v`: %w", status, err)
	}
	if err := redis.Client.Set(ctx, key, string(body), 0).Err(); err != nil {
		return err
	}
	return nil
}

func timeTextToMillis(text string) uint64 {
	var millis uint64
	if text != "" {
		t, err := time.Parse(TIME_LAYOUT, text)
		if err != nil {
			log.Warnf("failed to parse time format `%s`: %s", text, err)
		}
		millis = uint64(t.UnixMilli())
	}
	return millis
}

func Process(alert PostableAlert) {
	ctx := context.Background()
	hash := utils.ComputeHash(alert.Labels)
	key := fmt.Sprintf("alertmanager/%s", hash)
	startMillis := timeTextToMillis(alert.StartsAt)
	endMillis := timeTextToMillis(alert.EndsAt)

	status, err := getAlertStatus(ctx, hash)
	if err != nil {
		log.Warnf("failed to get alert status from redis for hash=%s: %s", hash, err)
		return
	}

	// New alert
	if status == nil {
		item := &api.Alert{
			Hash: hash,
			Source: api.Source{Kind: ALERTMANAGER_KIND, Name: config.InstanceName},
			StartsAt: startMillis,
			Labels: alert.Labels,
			Message: alert.Annotations["message"],
			Summary: alert.Annotations["summary"],
		}
		id, err := opensearch.StoreAlert(ctx, item)
		if err != nil {
			log.Warnf("failed to insert alert: %s", err)
			return
		}
		status = &api.AlertStatus{
			ID: id,
			LastCheck: startMillis,
			NextCheck: endMillis,
		}
		if err := setAlertStatus(ctx, key, status); err != nil {
			log.Warnf("failed to set alert status")
			return
		}
		return
	}

	// Update status
	newStatus := &api.AlertStatus{
		ID: status.ID,
		LastCheck: startMillis,
		NextCheck: endMillis,
	}
	if err := setAlertStatus(ctx, key, newStatus); err != nil {
		log.Warnf("failed to set alert status")
		return
	}
}
