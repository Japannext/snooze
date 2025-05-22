package tester

import (
	"os"
	"time"

	"github.com/japannext/snooze/internal/syslog"
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
		{Timestamp: today("06:24:43"), AppName: "sshd", ProcID: "", Host: "prod-host01", Severity: syslog.LOG_ERR, Facility: syslog.LOG_AUTH, Msg: "error: kex_exchange_identification: Connection closed by remote host"},
		{Timestamp: today("06:24:45"), AppName: "sshd", ProcID: "", Host: "prod-host02", Severity: syslog.LOG_ERR, Facility: syslog.LOG_AUTH, Msg: "error: kex_exchange_identification: Connection closed by remote host"},
		{Timestamp: today("06:36:40"), AppName: "sshd", ProcID: "", Host: "prod-host01", Severity: syslog.LOG_ERR, Facility: syslog.LOG_AUTH, Msg: "error: connect_to 10.1.2.3 port 443: failed"},
		{Timestamp: today("06:37:23"), AppName: "sshd", ProcID: "", Host: "prod-host01", Severity: syslog.LOG_ERR, Facility: syslog.LOG_AUTH, Msg: "error: PAM: User not known to the underlying authentication module for illegal user john.doe from workstation01.example.com"},
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
