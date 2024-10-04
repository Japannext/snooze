package alertmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/utils"
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

const ALERTMANAGER_KIND = "alertmanager"

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

const TIME_LAYOUT = "2006-01-02T15:04:05.999999999Z"

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

func hasKeys(labels map[string]string, keys ...string) bool {
	for _, key := range keys {
		if value, ok := labels[key]; !ok || value == "" {
			return false
		}
	}
	return true
}

const IDENTITY_PREFIX = "identity__"

func pop(m map[string]string, key string) string {
	v, ok := m[key]
	if ok {
		delete(m, key)
	}
	return v
}

func parseIdentity(labels map[string]string) map[string]string {
	identity := map[string]string{}
	for key, value := range labels {
		if strings.HasPrefix(key, IDENTITY_PREFIX) {
			k := strings.TrimPrefix(key, IDENTITY_PREFIX)
			identity[k] = value
		}
	}
	return identity
}

func parseSeverity(labels map[string]string) (string, int32) {
	text, ok := labels["severity__text"]
	if !ok {
		text = labels["severity"]
	}

	var number int32
	if n, ok := labels["severity__number"]; ok {
		i, err := strconv.Atoi(n)
		if err != nil {
			log.Warnf("in label `severity__number`, `%s` is not a number", n)
		}
		number = int32(i)
	}

	if number == 0 {
		number = utils.GuessSeverityNumber(text)
	}

	return text, number
}

func Process(alert PostableAlert) {
	ctx := context.Background()
	labels := alert.Labels
	alertName := pop(labels, "alertname")
	alertGroup := pop(labels, "alertgroup")
	hash := utils.ComputeHash(labels)
	key := fmt.Sprintf("alertmanager/%s/%s/%s", alertGroup, alertName, hash)
	startMillis := timeTextToMillis(alert.StartsAt)
	endMillis := timeTextToMillis(alert.EndsAt)

	status, err := getAlertStatus(ctx, key)
	if err != nil {
		log.Warnf("failed to get alert status from redis for hash=%s: %s", hash, err)
		return
	}

	// New alert
	if status == nil {
		annotations := alert.Annotations
		description := pop(annotations, "description")
		summary := pop(annotations, "summary")
		text, number := parseSeverity(labels)

		item := &api.Alert{
			Hash: hash,
			Source: api.Source{Kind: ALERTMANAGER_KIND, Name: config.InstanceName},
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
		id, err := opensearch.Store(ctx, api.ALERT_INDEX, item)
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
		ingestedAlerts.WithLabelValues(ALERTMANAGER_KIND, config.InstanceName).Inc()
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
	updatedAlerts.WithLabelValues(ALERTMANAGER_KIND, config.InstanceName).Inc()
}
