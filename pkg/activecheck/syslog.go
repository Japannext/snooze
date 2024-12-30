package activecheck

import (
	"fmt"
	"time"

	"github.com/japannext/snooze/pkg/common/syslog"
)

type SyslogConfig struct {
	Address string `validate:"required" yaml:"address"`
	Port    int    `yaml:"port"`
	// Valid values: rfc5424, rfc3164
	Format string `validate:"oneof=rfc5424 rfc3164" yaml:"format"`
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

func (probe *SyslogProbe) Fire(_ *Check, url string) error {
	client, err := syslog.NewClient(probe.addr, probe.format, time.Second)
	if err != nil {
		return fmt.Errorf("error initializing client: %w", err)
	}

	item := syslog.Log{
		Timestamp: time.Now(),
		AppName:   "snooze.activecheck",
		ProcID:    "",
		Host:      "",
		Severity:  syslog.LOG_DEBUG,
		Facility:  syslog.LOG_LOCAL0,
		Msg:       url,
	}
	if err := client.Send(item); err != nil {
		return fmt.Errorf("error sending log: %w", err)
	}

	return nil
}
