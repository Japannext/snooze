package syslog

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

type SyslogServer struct {
	srv *syslog.Server
}

type Handler struct {
}
func (h *Handler) Handle(record format.LogParts, msgLength int64, err error) {
	log.Debugf("received log: %s", record)
	ctx := context.TODO()
	ctx, span := tracer.Start(ctx, "syslog")
	defer span.End()
	if err != nil {
		log.Warnf("error handling log: %s", err)
		return
	}
	item := parseLog(ctx, record)
	if err := processQ.Publish(ctx, item); err != nil {
		log.Warnf("failed to publish log: %+v", err)
		return
	}
	ingestedLogs.WithLabelValues("syslog", config.InstanceName).Inc()
	log.Debugf("published log")
}

func NewSyslogServer() *SyslogServer {
	srv := syslog.NewServer()
	srv.SetFormat(syslog.RFC5424)
	srv.SetHandler(&Handler{})
	addr := fmt.Sprintf("%s:%d", config.ListenAddress, config.ListenPort)
	srv.ListenUDP(addr)
	srv.ListenTCP(addr)

	return &SyslogServer{srv}
}

func (s *SyslogServer) Run() error {
	s.srv.Boot()
	s.srv.Wait()

	return nil
}

func (s *SyslogServer) Stop() {
	s.srv.Kill()
}
