package syslog

import (
	"context"
	"time"
	"gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
	log "github.com/sirupsen/logrus"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

type SyslogServer struct {
	srv *syslog.Server
	ch syslog.LogPartsChannel
}

func NewSyslogServer() *SyslogServer {
	srv := syslog.NewServer()
	srv.SetFormat(syslog.RFC5424)
	ch := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(ch)
	srv.SetHandler(handler)
	srv.ListenUDP("0.0.0.0:1514")
	srv.ListenTCP("0.0.0.0:1514")

	return &SyslogServer{srv, ch}
}

var SEVERITY_TEXTS = []string{"emergency", "alert", "critical", "error", "warning", "notice", "informational", "debug"}
var SEVERITY_NUMBERS = []int32{21, 19, 18, 17, 13, 10, 9, 5}

var instanceName = ""

func parseLog(record format.LogParts) *api.Log {
	item := &api.Log{}
	item.Identity = make(map[string]string)
	item.Labels = make(map[string]string)

	item.TimestampMillis = uint64(record["timestamp"].(time.Time).UnixMilli())
	item.Source = api.Source{Kind: "syslog", Name: instanceName}
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

func (s *SyslogServer) Run() error {
	s.srv.Boot()

	go func (channel syslog.LogPartsChannel) {
		for record := range channel {
			ctx := context.Background()
			log.Debugf("Received log: %s", record)
			item := parseLog(record)

			if err := processChannel.Publish(ctx, item); err != nil {
				// TODO
				log.Warnf("Failed to publish log to process channel: %s", err)
				continue
			}
			log.Debug("Sent log to process channel")

		}
	}(s.ch)

	s.srv.Wait()

	return nil
}

func (s *SyslogServer) HandleStop() {
	close(s.ch)
	s.srv.Kill()
}
