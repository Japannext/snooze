package alertmanager

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/google/uuid"

	"github.com/japannext/snooze/pkg/common/utils"
	"github.com/japannext/snooze/pkg/models"
)

const ALERTMANAGER_KIND = "alertmanager"

func pop(m map[string]string, key string) string {
	v, ok := m[key]
	if ok {
		delete(m, key)
	}
	return v
}

func alertHandler(alert PostableAlert) {
	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "alertHandler")
	defer span.End()

	labels := alert.Labels
	var (
		alertName = pop(labels, "alertname")
		alertGroup = pop(labels, "alertgroup")
		hash = utils.ComputeHash(labels)
		key = fmt.Sprintf("alertmanager/%s/%s/%s", alertGroup, alertName, hash)
		startMillis = timeTextToMillis(alert.StartsAt)
		endMillis = timeTextToMillis(alert.EndsAt)
	)

	status, err := getAlertStatus(ctx, key)
	if err != nil {
		log.Warnf("failed to get alert status from redis for hash=%s: %s", hash, err)
		return
	}

	// New alert
	if status == nil {
		var (
			annotations = alert.Annotations
			description = pop(annotations, "description")
			summary = pop(annotations, "summary")
			text, number = parseSeverity(labels)
		)

		id := uuid.NewString()

		item := &models.Alert{
			ID: id,
			Hash: hash,
			Source: models.Source{Kind: ALERTMANAGER_KIND, Name: config.InstanceName},
			Identity: parseIdentity(labels),
			StartsAt: startMillis,
			AlertName: alertName,
			AlertGroup: alertGroup,
			SeverityText: text,
			SeverityNumber: number,
			Labels: alert.Labels,
			Description: description,
			Summary: summary,
		}
		if err := storeQ.Publish(ctx, item); err != nil {
			log.Warnf("failed to insert alert: %s", err)
			return
		}
		status = &models.AlertStatus{
			ID: id,
			LastCheck: startMillis,
			NextCheck: endMillis,
		}
		if err := setAlertStatus(ctx, key, status); err != nil {
			log.Warnf("failed to set alert status")
			return
		}
		ingestedAlerts.WithLabelValues(ALERTMANAGER_KIND, config.InstanceName).Inc()
		return
	}

	// Update status
	newStatus := &models.AlertStatus{
		ID: status.ID,
		LastCheck: startMillis,
		NextCheck: endMillis,
	}
	if err := setAlertStatus(ctx, key, newStatus); err != nil {
		log.Warnf("failed to set alert status")
		return
	}
	updatedAlerts.WithLabelValues(ALERTMANAGER_KIND, config.InstanceName).Inc()
}
