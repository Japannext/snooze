package syslog

import (
	"time"

	"gopkg.in/mcuadros/go-syslog.v2/format"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

var SEVERITY_TEXTS = []string{"emergency", "alert", "critical", "error", "warning", "notice", "informational", "debug"}
var SEVERITY_NUMBERS = []int32{21, 19, 18, 17, 13, 10, 9, 5}

const (
	SOURCE_KIND = "syslog"
)

type Parser struct {}

func NewParser() *Parser {
	return &Parser{}
}

func (parser *Parser) Run() error {
	for record := range receiveQueue {
		item := parseLog(record)
		publishQueue.Add(*item)
	}

	return nil
}

func (parser *Parser) Stop() {
}

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

	timestamp := record["timestamp"].(time.Time)
	observedTimestamp := time.Now()
	item.TimestampMillis = uint64(timestamp.UnixMilli())
	item.ObservedTimestampMillis = uint64(observedTimestamp.UnixMilli())
	if item.TimestampMillis == 0 {
		item.TimestampMillis = item.ObservedTimestampMillis
		emptyTimestamp.WithLabelValues(SOURCE_KIND, config.InstanceName).Inc()
	} else {
		delay := observedTimestamp.Sub(timestamp).Seconds()
		sourceDelay.WithLabelValues(SOURCE_KIND, config.InstanceName).Observe(delay)
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

	// Active-check
	if item.Identity["process"] == "snooze.activecheck" {
		item.Labels["activecheck.url"] = item.Message
	}

	return item
}
