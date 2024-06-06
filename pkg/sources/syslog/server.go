package syslog

import (
	"context"
	"gopkg.in/mcuadros/go-syslog.v2"
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

func (s *SyslogServer) Run() error {
	s.srv.Boot()

	go func (channel syslog.LogPartsChannel) {
		for record := range channel {
			ctx := context.Background()
			log.Debugf("Received log: %s", record)
			alert := &api.Alert{
				Source: api.Source{Kind: "syslog", Name: ""},
				Labels: map[string]string{
					"host.name": record["hostname"].(string),
					"process": record["app_name"].(string),
				},
				Body: map[string]string{
					"message": record["message"].(string),
				},
			}
			severity, found := record["severity"].(int)
			if found  && severity >= 0 && severity < 7 {
				alert.SeverityText = SEVERITY_TEXTS[severity]
				alert.SeverityNumber = SEVERITY_NUMBERS[severity]
			}
			if err := processChannel.Publish(ctx, alert); err != nil {
				// TODO
				log.Warnf("Failed to publish log to process channel: %s", err)
				continue
			}
			log.Debug("Sent alert to process channel")

		}
	}(s.ch)

	s.srv.Wait()

	return nil
}

func (s *SyslogServer) HandleStop() {
	close(s.ch)
	s.srv.Kill()
}
