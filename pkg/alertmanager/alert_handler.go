package alertmanager

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/japannext/snooze/pkg/common/opensearch/format"
	"github.com/japannext/snooze/pkg/common/utils"
	"github.com/japannext/snooze/pkg/models"
	log "github.com/sirupsen/logrus"
)

const ALERTMANAGER_KIND = "alertmanager"

func pop(m map[string]string, key string) string {
	v, ok := m[key]
	if ok {
		delete(m, key)
	}
	return v
}

func alertHandler(ctx context.Context, alert PostableAlert) {
	ctx, span := tracer.Start(ctx, "alertHandler")
	defer span.End()

	labels := alert.Labels
	var (
		alertName  = pop(labels, "alertname")
		alertGroup = pop(labels, "alertgroup")
		hash       = utils.ComputeHash(labels)
		key        = fmt.Sprintf("alertmanager/%s/%s/%s", alertGroup, alertName, hash)
		startTime  = parseTime(alert.StartsAt)
		endTime    = parseTime(alert.EndsAt)
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
			summary     = pop(annotations, "summary")
		)

		id := uuid.NewString()

		item := &models.Alert{}

		item.SetID(id)
		item.Hash = hash
		item.Source.Kind = ALERTMANAGER_KIND
		item.Source.Name = config.InstanceName
		item.Identity = parseIdentity(labels)
		item.StartAt = startTime
		item.AlertName = alertName
		item.AlertGroup = alertGroup
		item.SeverityText, item.SeverityNumber = parseSeverity(labels)
		item.Labels = alert.Labels
		item.Description = description
		item.Summary = summary

		err := storeQ.PublishData(ctx, &format.Create{
			Index: models.ALERT_INDEX,
			Item:  item,
		})
		if err != nil {
			log.Warnf("failed to insert alert: %s", err)
			return
		}
		status = &models.AlertStatus{
			ID:        id,
			LastCheck: startTime,
			NextCheck: endTime,
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
		ID:        status.ID,
		LastCheck: startTime,
		NextCheck: endTime,
	}
	if err := setAlertStatus(ctx, key, newStatus); err != nil {
		log.Warnf("failed to set alert status")
		return
	}
	updatedAlerts.WithLabelValues(ALERTMANAGER_KIND, config.InstanceName).Inc()
}
