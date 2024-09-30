package syslog

import (
	"time"

	"gopkg.in/mcuadros/go-syslog.v2/format"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

const (
	SOURCE_KIND = "syslog"
)

func extract(record format.LogParts, key string) (string, bool) {
	text, ok := record[key].(string)
	if !ok || text == "-" || text == "" {
		return "", false
	}
	return text, true
}

func parseLog(record format.LogParts) *api.Log {
	item := &api.Log{}
	item.Identity = make(map[string]string)
	item.Labels = make(map[string]string)

	item.TimestampMillis = uint64(record["timestamp"].(time.Time).UnixMilli())
	item.ObservedTimestampMillis = uint64(time.Now().UnixMilli())
	if item.TimestampMillis == 0 {
		item.TimestampMillis = item.ObservedTimestampMillis
	}
	item.Source = api.Source{Kind: SOURCE_KIND, Name: config.InstanceName}

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

	return item
}
