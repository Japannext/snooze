package activecheck

import (
	"fmt"
	"time"

	"github.com/japannext/snooze/pkg/common/syslog"
)

type SyslogConfig struct {
	Address string `yaml:"address" validate:"required"`
	Port    int    `yaml:"port"`
	// Valid values: rfc5424, rfc3164
	Format string `yaml:"format" validate:"oneof=rfc5424 rfc3164"`
}

type SyslogProbe struct {
	name   string
	addr   string
	format string
}

func NewSyslogProbe(name string, cfg *SyslogConfig) Probe {
	if cfg.Format == "" {
		cfg.Format = "rfc5424"
	}
	if cfg.Port == 0 {
		cfg.Port = 514
	}
	return &SyslogProbe{
		name:   name,
		addr:   fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
		format: cfg.Format,
	}
}

func (probe *SyslogProbe) Fire(check *Check, url string) error {
	client, err := syslog.NewClient(probe.addr, probe.format, time.Second)
	if err != nil {
		return err
	}

	item := syslog.Log{
		Timestamp: time.Now(),
		AppName:   "snooze.activecheck",
		ProcId:    "",
		Host:      "",
		Severity:  syslog.LOG_DEBUG,
		Facility:  syslog.LOG_LOCAL0,
		Msg:       url,
	}
	if err := client.Send(item); err != nil {
		return err
	}
	return nil
}
