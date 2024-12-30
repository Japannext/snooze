package syslog

import (
	"context"
	"time"

	"github.com/japannext/snooze/pkg/common/tracing"
	"github.com/japannext/snooze/pkg/models"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

var (
	SEVERITY_TEXTS   = []string{"emergency", "alert", "critical", "error", "warning", "notice", "informational", "debug"}
	SEVERITY_NUMBERS = []int32{21, 19, 18, 17, 13, 10, 9, 5}
)

var FACILITY_TEXTS = []string{
	"kern", "user", "mail", "daemon", "auth", "syslog", "lpr", "news",
	"uupc", "cron", "authpriv", "ftp", "ntp", "security", "console", "solaris-cron",
	"local0", "local1", "local2", "local3", "local4", "local5", "local6", "local7",
}

const (
	SOURCE_KIND = "syslog"
)

func parseLog(ctx context.Context, record format.LogParts) *models.Log {
	ctx, span := tracer.Start(ctx, "parseLog")
	defer span.End()

	item := &models.Log{}
	item.Identity = make(map[string]string)
	item.Labels = make(map[string]string)
	item.TraceID = tracing.GetTraceID(ctx)

	// Time
	timestamp := record["timestamp"].(time.Time)
	observed := models.TimeNow()
	item.ActualTime = models.Time{Time: timestamp}
	item.ObservedTime = observed
	if item.ActualTime.IsZero() {
		emptyTimestamp.WithLabelValues(SOURCE_KIND, config.InstanceName).Inc()
	} else {
		delay := observed.Sub(timestamp).Seconds()
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

	facilityNumber := record["facility"].(int)
	if 0 <= facilityNumber && facilityNumber < len(FACILITY_TEXTS) {
		item.Labels["syslog.facility"] = FACILITY_TEXTS[facilityNumber]
	}

	if value, ok := extract(record, "proc_id"); ok {
		item.Labels["syslog.proc_id"] = value
	}
	if value, ok := extract(record, "msg_id"); ok {
		item.Labels["syslog.msg_id"] = value
	}

	item.Message = record["message"].(string)

	severity, found := record["severity"].(int)
	if found && severity >= 0 && severity < 7 {
		item.SeverityText = SEVERITY_TEXTS[severity]
		item.SeverityNumber = SEVERITY_NUMBERS[severity]
	}

	setSpanLog(span, record)

	// Active-check
	if item.Identity["process"] == "snooze.activecheck" {
		item.ActiveCheckURL = item.Message
	}

	return item
}

const TRACE_TIME_FORMAT = "2006-01-02T15:04:05 -07:00"

func setSpanLog(span trace.Span, record format.LogParts) {
	tracing.SetInt(span, "syslog.priority", record["priority"].(int))
	tracing.SetInt(span, "syslog.facility", record["facility"].(int))
	tracing.SetInt(span, "syslog.severity", record["severity"].(int))
	tracing.SetInt(span, "syslog.version", record["version"].(int))
	tracing.SetTime(span, "syslog.timestamp", record["timestamp"].(time.Time))
	tracing.SetString(span, "syslog.hostname", record["hostname"].(string))
	tracing.SetString(span, "syslog.app_name", record["app_name"].(string))
	tracing.SetString(span, "syslog.proc_id", record["proc_id"].(string))
	tracing.SetString(span, "syslog.msg_id", record["msg_id"].(string))
	tracing.SetString(span, "syslog.structured_data", record["structured_data"].(string))
	tracing.SetString(span, "syslog.message", record["message"].(string))
}

func extract(record format.LogParts, key string) (string, bool) {
	text, ok := record[key].(string)
	if !ok || text == "-" || text == "" {
		return "", false
	}
	return text, true
}
