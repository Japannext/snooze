package tester

import (
	"os"
	"time"

	"github.com/japannext/snooze/pkg/common/syslog"
	log "github.com/sirupsen/logrus"
)

var todayDate = time.Now().Format("2006-01-02")

func today(t string) time.Time {
	dt, err := time.Parse("2006-01-02T15:04:05", todayDate+"T"+t)
	if err != nil {
		log.Fatal(err)
	}

	return dt
}

func syslogSample() []syslog.Log {
	return []syslog.Log{
		{today("06:24:43"), "sshd", "", "prod-host01", syslog.LOG_ERR, syslog.LOG_AUTH, "error: kex_exchange_identification: Connection closed by remote host"},
		{today("06:24:45"), "sshd", "", "prod-host02", syslog.LOG_ERR, syslog.LOG_AUTH, "error: kex_exchange_identification: Connection closed by remote host"},
		{today("06:36:40"), "sshd", "", "prod-host01", syslog.LOG_ERR, syslog.LOG_AUTH, "error: connect_to 10.1.2.3 port 443: failed"},
		{today("06:37:23"), "sshd", "", "prod-host01", syslog.LOG_ERR, syslog.LOG_AUTH, "error: PAM: User not known to the underlying authentication module for illegal user john.doe from workstation01.example.com"},
	}
}

func runSyslogSamples() error {
	serverAddr := os.Getenv("SYSLOG_SERVER_ADDR")

	client, err := syslog.NewClient(serverAddr, "rfc5424", time.Second)
	if err != nil {
		return err
	}

	if err := client.Send(syslogSample()...); err != nil {
		return err
	}

	return nil
}
