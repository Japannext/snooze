package syslog

import (
	"context"
	"time"

	"gopkg.in/mcuadros/go-syslog.v2/format"

	"github.com/japannext/snooze/pkg/models"
)

var SEVERITY_TEXTS = []string{"emergency", "alert", "critical", "error", "warning", "notice", "informational", "debug"}
var SEVERITY_NUMBERS = []int32{21, 19, 18, 17, 13, 10, 9, 5}

const (
	SOURCE_KIND = "syslog"
)

func parseLog(ctx context.Context, record format.LogParts) *models.Log {
	ctx, span := tracer.Start(ctx, "parseLog")
	defer span.End()

	item := &models.Log{}
	item.Identity = make(map[string]string)
	item.Labels = make(map[string]string)

	timestamp := record["timestamp"].(time.Time)
	observedTimestamp := time.Now()
	item.Timestamp.Actual = uint64(timestamp.UnixMilli())
	item.Timestamp.Observed = uint64(observedTimestamp.UnixMilli())
	if item.Timestamp.Actual == 0 {
		emptyTimestamp.WithLabelValues(SOURCE_KIND, config.InstanceName).Inc()
	} else {
		delay := observedTimestamp.Sub(timestamp).Seconds()
		sourceDelay.WithLabelValues(SOURCE_KIND, config.InstanceName).Observe(delay)
	}

	item.Source = models.Source{Kind: SOURCE_KIND, Name: config.InstanceName}

	clientIP := record["client"].(string)

	// Identity
	if value, ok := extract(record, "hostname"); ok {
		item.Identity["kind"] = "host"
		item.Identity["host"] = value
	} else {
		item.Identity["kind"] = "ip"
		item.Identity["ip"] = clientIP
	}
	if value, ok := extract(record, "app_name"); ok {
		item.Identity["process"] = value
	}

	// Labels
	item.Labels["syslog.client_ip"] = clientIP
	if value, ok := extract(record, "tls_peer"); ok {
		item.Labels["syslog.tls_peer"] = value
	}

	if value, ok := extract(record, "proc_id"); ok {
		item.Labels["syslog.proc_id"] = value
	}
	if value, ok := extract(record, "msg_id"); ok {
		item.Labels["syslog.msg_id"] = value
	}

	item.Message = record["message"].(string)

	severity, found := record["severity"].(int)
	if found  && severity >= 0 && severity < 7 {
		item.SeverityText = SEVERITY_TEXTS[severity]
		item.SeverityNumber = SEVERITY_NUMBERS[severity]
	}

	// Active-check
	if item.Identity["process"] == "snooze.activecheck" {
		item.ActiveCheckURL = item.Message
	}

	return item
}

func extract(record format.LogParts, key string) (string, bool) {
	text, ok := record[key].(string)
	if !ok || text == "-" || text == "" {
		return "", false
	}
	return text, true
}
