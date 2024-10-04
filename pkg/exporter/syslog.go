package exporter

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/common/syslog"
)

type SyslogRelay struct {
	name string
	addr string
	format string
}

func NewSyslogRelay(name string, cfg *SyslogConfig) ProbeHandler {
	if cfg.Format == "" {
		cfg.Format = "rfc5424"
	}
	if cfg.Port == 0 {
		cfg.Port = 514
	}
	return &SyslogRelay{
		name: name,
		addr: fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
		format: cfg.Format,
	}
}

func (relay *SyslogRelay) ServeHTTP(c *gin.Context) {
	var (
		syslogUp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "snooze_activecheck",
			Name: "syslog_up",
			Help: "Active check that verify syslog messages are processed",
		}, []string{"relay"})
		syslogDelay = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "snooze_activecheck",
			Name: "syslog_delay",
			Help: "Time the message took from sending to delivery",
			Buckets: prometheus.ExponentialBuckets(0.1, 2, 8),
		}, []string{"relay"})
	)
	reg := prometheus.NewRegistry()
	reg.MustRegister(syslogUp, syslogDelay)

	key := uuid.NewString()
	log.Debugf("[relay=%s, key=%s] Checking syslog", relay.name, key)

	if err := waiter.Prepare(key); err != nil {
		log.Error(err)
		return
	}
	defer waiter.Cleanup(key)

	// Send a syslog message
	client, err := syslog.NewClient(relay.addr, relay.format, time.Second)
	if err != nil {
		log.Error(err)
		return
	}

	callbackURL := fmt.Sprintf("http://%s:%d/webhook/%s", config.CallbackAddress, config.CallbackPort, key)
	log.Debugf("[relay=%s, key=%s] callback URL: %s", relay.name, key, callbackURL)
	item := syslog.Log{
		Timestamp: time.Now(),
		AppName: "snooze.activecheck",
		ProcId: "",
		Host: "",
		Severity: syslog.LOG_DEBUG,
		Facility: syslog.LOG_LOCAL0,
		Msg: callbackURL,
	}
	log.Debugf("[relay=%s, key=%s] sending syslog message (format=%s)", relay.name, key, relay.format)
	if err := client.Send(item); err != nil {
		log.Error(err)
		return
	}

	// Waiting for answer from snooze-process
	callback, err := waiter.Wait(key)
	if err != nil {
		log.Warnf("[relay=%s, key=%s] %s", relay.name, key, err)
		return
	}

	log.Debugf("[relay=%s, key=%s] received callback", relay.name, key)
	// Update the metrics
	if callback.Error != "" {
		syslogUp.WithLabelValues(relay.name).Set(0)
	} else {
		syslogUp.WithLabelValues(relay.name).Set(1)
	}
	syslogDelay.WithLabelValues(relay.name).Observe(float64(callback.DelayMillis / 1000))

	promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(c.Writer, c.Request)
}
