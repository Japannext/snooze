package syslog

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

var receiveQueue = make(syslog.LogPartsChannel)

type SyslogServer struct {
	srv *syslog.Server
}

type Handler struct {}
func (h *Handler) Handle(record format.LogParts, msgLength int64, err error) {
	ctx := context.TODO()
	if err != nil {
		log.Warnf("error handling log: %s", err)
		return
	}
	item := parseLog(ctx, record)
	processQ.Publish(ctx, item)
}

func NewSyslogServer() *SyslogServer {
	srv := syslog.NewServer()
	srv.SetFormat(syslog.RFC5424)
	handler := syslog.NewChannelHandler(receiveQueue)
	srv.SetHandler(handler)
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
	close(receiveQueue)
	s.srv.Kill()
}
