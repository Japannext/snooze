package syslog

import (
	"time"

	"gopkg.in/mcuadros/go-syslog.v2/format"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

const (
	SOURCE_KIND = "syslog"
)

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
	item.Identity["kind"] = "host"
	item.Identity["hostname"] = record["hostname"].(string)
	item.Identity["process"] = record["app_name"].(string)

	item.Labels["client"] = record["client"].(string)
	if tlsPeer := record["tls_peer"].(string); tlsPeer != "" {
		item.Labels["tls_peer"] = tlsPeer
	}

	item.Labels["proc_id"] = record["proc_id"].(string)
	item.Labels["msg_id"] = record["msg_id"].(string)

	item.Message = record["message"].(string)

	severity, found := record["severity"].(int)
	if found  && severity >= 0 && severity < 7 {
		item.SeverityText = SEVERITY_TEXTS[severity]
		item.SeverityNumber = SEVERITY_NUMBERS[severity]
	}

	return item
}
