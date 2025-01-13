package alertmanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/models"
)

// Insert a status into the database
func insertStatus(ctx context.Context, key string, status *models.AlertStatus) error {

	return nil
}

func getAlertStatus(ctx context.Context, key string) (*models.AlertStatus, bool, error) {
	body, err := redis.Client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return &models.AlertStatus{}, false, nil
	}

	if err != nil {
		return &models.AlertStatus{}, false, fmt.Errorf("failed to get redis key '%s': %w", key, err)
	}

	var status *models.AlertStatus

	if err := json.Unmarshal(body, &status); err != nil {
		return &models.AlertStatus{}, false, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return status, true, nil
}

func setAlertStatus(ctx context.Context, key string, status *models.AlertStatus) error {
	body, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal alert status `%+v`: %w", status, err)
	}

	if err := redis.Client.Set(ctx, key, string(body), 0).Err(); err != nil {
		return fmt.Errorf("failed to set alert status: %w", err)
	}
	return nil
}
