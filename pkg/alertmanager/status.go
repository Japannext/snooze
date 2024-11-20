package alertmanager

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/redis"
)

func getAlertStatus(ctx context.Context, key string) (*models.AlertStatus, error) {
	body, err := redis.Client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if err == redis.Nil {
		return nil, nil
	}
	var status *models.AlertStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, err
	}
	return status, nil
}

func setAlertStatus(ctx context.Context, key string, status *models.AlertStatus) error {
	body, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal alert status `%+v`: %w", status, err)
	}
	if err := redis.Client.Set(ctx, key, string(body), 0).Err(); err != nil {
		return err
	}
	return nil
}
