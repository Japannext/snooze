package exporter

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/japannext/snooze/pkg/common/syslog"
)

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
		Buckets: []float64{0.1, 0.5, 1.0, 5.0, 10.0},
	}, []string{"relay"})
)

func initSyslogMetrics() {
	prometheus.MustRegister(
		syslogUp,
		syslogDelay,
	)
}

// Configuration
type SyslogRelay struct {
	Name string `yaml:"name" validate:"required"`
	Address string `yaml:"address" validate:"required"`
	Port int `yaml:"port"`
	// Valid values: rfc5424, rfc3164
	Format string `yaml:"format" validate:"oneof=rfc5424 rfc3164"`
}

func (relay *SyslogRelay) FillDefaults() {
	if relay.Format == "" {
		relay.Format = "rfc5424"
	}
	if relay.Port == 0 {
		relay.Port = 514
	}
}

func (relay *SyslogRelay) Check() error {
	key := uuid.NewString()

	if err := waiter.Prepare(key); err != nil {
		return err
	}
	defer waiter.Cleanup(key)

	// Send a syslog message
	client, err := syslog.NewClient(fmt.Sprintf("%s:%d", relay.Address, relay.Port), relay.Format, time.Second)
	if err != nil {
		return err
	}

	callbackURL := fmt.Sprintf("http://%s:%d/webhook/%s", config.CallbackAddress, config.CallbackPort, key)
	item := syslog.Log{
		Timestamp: time.Now(),
		AppName: "snooze.activecheck",
		ProcId: "",
		Host: "",
		Severity: syslog.LOG_DEBUG,
		Facility: syslog.LOG_LOCAL0,
		Msg: callbackURL,
	}
	switch relay.Format {
		case "rfc5424":
		default:
			log.Fatalf("unsupported format: %s", relay.Format)
	}
	if err := client.Send(item); err != nil {
		return err
	}

	// Waiting for answer from snooze-process
	callback, err := waiter.Wait(key)
	if err != nil {
		syslogUp.WithLabelValues(relay.Name).Set(0)
	}

	// Update the metrics
	if callback.Error != "" {
		syslogUp.WithLabelValues(relay.Name).Set(0)
	} else {
		syslogUp.WithLabelValues(relay.Name).Set(1)
	}
	syslogDelay.WithLabelValues(relay.Name).Observe(float64(callback.DelayMillis / 1000))

	return nil
}
