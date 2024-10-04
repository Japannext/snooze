package syslog

import (
	"fmt"

	"gopkg.in/mcuadros/go-syslog.v2"
)

var receiveQueue = make(syslog.LogPartsChannel)

type SyslogServer struct {
	srv *syslog.Server
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
