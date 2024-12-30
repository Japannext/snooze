package syslog

import (
	"bytes"
	"fmt"
	"net"
	"time"
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
	Timestamp time.Time
	AppName   string
	ProcId    string
	Host      string
	Severity  int
	Facility  int
	Msg       string
}

func (item *Log) toRFC3164() []byte {
	var buf bytes.Buffer
	pri := item.Severity | item.Facility
	buf.WriteString(fmt.Sprintf("<%d> ", pri))
	buf.WriteString(item.Timestamp.Format("Jan 2 15:04:05"))
	buf.WriteString(" ")
	buf.WriteString(item.Host)
	buf.WriteString(" ")
	buf.WriteString(item.AppName)
	if item.ProcId != "" {
		buf.WriteString("[")
		buf.WriteString(item.ProcId)
		buf.WriteString("]")
	}
	buf.WriteString(": ")
	buf.WriteString(item.Msg)
	buf.WriteString("\n")

	return buf.Bytes()
}

func (item *Log) toRFC5424() []byte {
	var buf bytes.Buffer
	pri := item.Severity | item.Facility
	buf.WriteString(fmt.Sprintf("<%d>1 ", pri))
	buf.WriteString(item.Timestamp.Format("2006-01-02T15:04:05.999Z"))
	buf.WriteString(" ")
	buf.WriteString(item.Host)
	buf.WriteString(" ")
	buf.WriteString(item.AppName)
	buf.WriteString(" ")
	if item.ProcId != "" {
		buf.WriteString(item.ProcId)
	} else {
		buf.WriteString("-")
	}
	buf.WriteString(" ")
	buf.WriteString("- - ")
	buf.WriteString(item.Msg)
	buf.WriteString("\n")

	return buf.Bytes()
}

type Client struct {
	net.Conn
	Format string
}

func NewClient(addr string, format string, timeout time.Duration) (*Client, error) {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}
	return &Client{conn, format}, nil
}

func (client *Client) Send(items ...Log) error {
	for _, item := range items {
		var data []byte
		switch client.Format {
		case "rfc5424":
			data = item.toRFC5424()
		case "rfc3164":
			data = item.toRFC3164()
		}
		client.Write(data)
	}
	return nil
}
