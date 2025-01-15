package alertmanager

import (
	"context"
	"fmt"
	"time"

	"github.com/japannext/snooze/pkg/common/opensearch/format"
	"github.com/japannext/snooze/pkg/common/utils"
	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/alertmanager/status"
	"github.com/japannext/snooze/pkg/common/tracing"
	log "github.com/sirupsen/logrus"
)

const (
	ALERTMANAGER_KIND = "alertmanager"
	opensearchIDByteSize = 30
)

func alertHandler(ctx context.Context, alert PostableAlert) {
	ctx, span := tracer.Start(ctx, "alertHandler")
	defer span.End()

	key := getKey(alert)

	// Check if the alert is already active (in Redis).
	alertStatus, active, err := isAlertActive(ctx, key)
	if err != nil {
		log.Warnf("failed to insert alert: %s", err)

		return
	}

	// Update the Redis value if the alert is active.
	if active {
		if err := updateActiveAlert(ctx, alertStatus, alert); err != nil {
			log.Warnf("failed to update alert: %s", err)
		}

		return
	}

	newActiveAlert(ctx, key, alert)
}

func getKey(alert PostableAlert) string {
	hash := utils.ComputeHash(alert.Labels)

	return fmt.Sprintf("alertmanager/%s/%s/%s", alert.Labels["alertgroup"], alert.Labels["alertname"], hash)
}

func isAlertActive(ctx context.Context, key string) (*status.AlertStatus, bool, error) {
	ctx, span := tracer.Start(ctx, "isAlertActive")
	defer span.End()

	alertStatus, found, err := status.Get(ctx, key)
	if err != nil {
		return &status.AlertStatus{}, false, fmt.Errorf("isAlertActive failed for key %s: %w", key, err)
	}

	tracing.SetBool(span, "found", found)
	tracing.SetTime(span, "status.nextCheck", alertStatus.NextCheck.Time)

	if !found {
		return &status.AlertStatus{}, false, nil
	}

	// In redis, but not cleaned-up yet
	if alertStatus.NextCheck.Before(time.Now()) {
		return &status.AlertStatus{}, false, nil
	}

	return alertStatus, true, nil
}

func updateActiveAlert(ctx context.Context, alertStatus *status.AlertStatus, alert PostableAlert) error {
	ctx, span := tracer.Start(ctx, "updateActiveAlert")
	defer span.End()

	alertStatus.LastCheck = parseTime(alert.StartsAt)
	alertStatus.NextCheck = parseTime(alert.EndsAt)

	err := status.Set(ctx, alertStatus)
	if err != nil {
		return fmt.Errorf("failed to update alert status: %w", err)
	}

	err = storeQ.PublishData(ctx, &format.Update{
		Index: models.ActiveAlertIndex,
		ID: alertStatus.ID,
		Doc: struct{LastHit models.Time `json:"lastHit"`}{
			LastHit: models.TimeNow(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to queue active alert update: %w", err)
	}

	updatedAlerts.WithLabelValues(ALERTMANAGER_KIND, config.InstanceName).Inc()

	return nil
}

func newActiveAlert(ctx context.Context, key string, alert PostableAlert) {
	ctx, span := tracer.Start(ctx, "newActiveAlert")
	defer span.End()

	activeAlert := convertAlert(alert)
	id := utils.RandomURLSafeBase64(opensearchIDByteSize)

	activeAlert.TraceID = tracing.GetTraceID(ctx)

	// Add it to opensearch
	err := storeQ.PublishData(ctx, &format.Index{
		Index: models.ActiveAlertIndex,
		ID: id,
		Item: activeAlert,
	})
	if err != nil {
		log.Errorf("failed to queue active alert: %s", err)

		return
	}

	// Add alert status (redis)
	err = status.Set(ctx, &status.AlertStatus{
		ID: id,
		Key: key,
		LastCheck: parseTime(alert.StartsAt),
		NextCheck: parseTime(alert.EndsAt),
	})
	if err != nil {
		log.Errorf("failed to create alert status: %s", err)

		return
	}

	ingestedAlerts.WithLabelValues(ALERTMANAGER_KIND, config.InstanceName).Inc()
}
