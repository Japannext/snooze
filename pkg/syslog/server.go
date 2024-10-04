package syslog

import (
	"fmt"
	"time"

	"gopkg.in/mcuadros/go-syslog.v2"
	log "github.com/sirupsen/logrus"
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
	addr := fmt.Sprintf("%s:%d", config.ListenAddress, config.ListenPort)
	srv.ListenUDP(addr)
	srv.ListenTCP(addr)

	return &SyslogServer{srv, ch}
}

var SEVERITY_TEXTS = []string{"emergency", "alert", "critical", "error", "warning", "notice", "informational", "debug"}
var SEVERITY_NUMBERS = []int32{21, 19, 18, 17, 13, 10, 9, 5}

func (s *SyslogServer) Run() error {
	s.srv.Boot()

	go func (channel syslog.LogPartsChannel) {
		for record := range channel {
			start := time.Now()
			log.Debugf("Received log: %s", record)
			item := parseLog(record)

			if err := producer.Publish(item); err != nil {
				// TODO
				log.Warnf("Failed to publish log to process channel: %s", err)
				continue
			}
			log.Debug("Sent log to process channel")

			processMetrics(start, item)
		}
	}(s.ch)

	s.srv.Wait()

	return nil
}

func (s *SyslogServer) Stop() {
	close(s.ch)
	s.srv.Kill()
}
