package samples

import (
	"bytes"
	"os"
	"net"
	"time"
	"fmt"

	log "github.com/sirupsen/logrus"
)

const (
	// Severity.

	// From /usr/include/sys/syslog.h.
	// These are the same on Linux, BSD, and OS X.
	LOG_EMERG int = iota
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
)

const (
	// Facility.

	// From /usr/include/sys/syslog.h.
	// These are the same up to LOG_FTP on Linux, BSD, and OS X.
	LOG_KERN int = iota << 3
	LOG_USER
	LOG_MAIL
	LOG_DAEMON
	LOG_AUTH
	LOG_SYSLOG
	LOG_LPR
	LOG_NEWS
	LOG_UUCP
	LOG_CRON
	LOG_AUTHPRIV
	LOG_FTP
	_ // unused
	_ // unused
	_ // unused
	_ // unused
	LOG_LOCAL0
	LOG_LOCAL1
	LOG_LOCAL2
	LOG_LOCAL3
	LOG_LOCAL4
	LOG_LOCAL5
	LOG_LOCAL6
	LOG_LOCAL7
)


var todayDate = time.Now().Format("2006-01-02")

type Log struct {
	timestamp time.Time
	appName string
	procId string
	host string
	severity int
	facility int
	msg string
}

func (item *Log) toRFC5424() []byte {
	var buf bytes.Buffer
	pri := item.severity | item.facility
	buf.WriteString(fmt.Sprintf("<%d>1 ", pri))
	buf.WriteString(item.timestamp.Format("2006-01-02T15:04:05.999Z"))
	buf.WriteString(" ")
	buf.WriteString(item.host)
	buf.WriteString(" ")
	buf.WriteString(item.appName)
	buf.WriteString(" ")
	if item.procId != "" {
		buf.WriteString(item.procId)
	} else {
		buf.WriteString("-")
	}
	buf.WriteString(" ")
	buf.WriteString("- - ")
	buf.WriteString(item.msg)
	buf.WriteString("\n")

	return buf.Bytes()
}

func today(t string) time.Time {
	dt, err := time.Parse("2006-01-02T15:04:05", todayDate + "T" + t)
	if err != nil {
		log.Fatal(err)
	}
	return dt
}

var syslogSample = []Log{
	{today("06:24:43"), "sshd", "", "prod-host01", LOG_ERR, LOG_AUTH, "error: kex_exchange_identification: Connection closed by remote host"},
	{today("06:24:45"), "sshd", "", "prod-host02", LOG_ERR, LOG_AUTH, "error: kex_exchange_identification: Connection closed by remote host"},
	{today("06:36:40"), "sshd", "", "prod-host01", LOG_ERR, LOG_AUTH, "error: connect_to 10.1.2.3 port 443: failed"},
	{today("06:37:23"), "sshd", "", "prod-host01", LOG_ERR, LOG_AUTH, "error: PAM: User not known to the underlying authentication module for illegal user john.doe from workstation01.example.com"},
}

func runSyslogSamples() error {

	serverAddr := os.Getenv("SYSLOG_SERVER_ADDR")

	conn, err := net.DialTimeout("tcp", serverAddr, time.Second)
	if err != nil {
		return err
	}
	for _, item := range syslogSample {
		data := item.toRFC5424()
		log.Infof("[SAMPLE] Sending: %s", data)
		conn.Write(data)
	}
	return nil
}
