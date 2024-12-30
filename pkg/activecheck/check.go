package activecheck

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/models"
)

func probeHandler(c *gin.Context) {
	name := c.Param("name")
	check, found := checks[name]
	if !found {
		// TODO
	}
	check.Handle(c)
}

type Check struct {
	Name         string              `yaml:"name" validate:"required"`
	Source       *SourceConfig       `yaml:"source"`
	Notification *NotificationConfig `yaml:"notification"`
	Timeout      *time.Duration      `yaml:"timeout"`

	internal struct {
		probe Probe
	}
}

type Probe interface {
	Fire(*Check, string) error
}

type SourceConfig struct {
	Syslog *SyslogConfig `yaml:"syslog"`
}

type NotificationConfig struct {
	ExpectedQueue   string `yaml:"expected_queue"`
	ExpectedProfile string `yaml:"expected_profile"`
}

func (check *Check) Load() {
	if check.Name == "" {
		log.Fatalf("no `name` for check")
	}
	if check.Timeout == nil {
		timeout := 30 * time.Second
		check.Timeout = &timeout
	}
	switch {
	case check.Source.Syslog != nil:
		check.internal.probe = NewSyslogProbe(check.Name, check.Source.Syslog)
	default:
		log.Fatalf("no `source` defined in configuration for check '%s'", check.Name)
	}
}

var (
	probeUp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "snooze",
		Name:      "activecheck_probe_up",
		Help:      "Active check that verify a probe is up",
	}, []string{"check"})
)

func (check *Check) Handle(c *gin.Context) {

	// New prometheus registry
	reg := prometheus.NewRegistry()
	reg.MustRegister(probeUp)

	key := uuid.NewString()
	url := fmt.Sprintf("http://%s:%d/webhook/%s", config.CallbackAddress, config.CallbackPort, key)
	if err := check.internal.probe.Fire(check, url); err != nil {
		// TODO
	}

	callback, err := waiter.Wait(key, *check.Timeout)
	if err != nil {
		// TODO
	}

	ok := check.CheckCallback(callback)
	if ok {
		probeUp.WithLabelValues(check.Name).Set(1)
	} else {
		probeUp.WithLabelValues(check.Name).Set(0)
	}

	promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(c.Writer, c.Request)
}

func callbackHandler(c *gin.Context) {
	key := c.Param("uid")
	var callback models.SourceActiveCheck
	c.BindJSON(&callback)

	waiter.Insert(key, callback)
}

func (check *Check) CheckCallback(callback models.SourceActiveCheck) bool {
	if callback.Error != "" {
		return false
	}

	return true
}
