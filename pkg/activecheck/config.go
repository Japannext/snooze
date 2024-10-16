package activecheck

import (
	"net"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var config *Config
var probes = map[string]ProbeHandler{}

type Config struct {
	CallbackAddress string `mapstructure:"CALLBACK_ADDRESS"`
	CallbackPort int `mapstructure:"CALLBACK_PORT"`
	SyslogAddress string `mapstructure:"SYSLOG_ADDRESS"`
	SyslogPort int `mapstructure:"SYSLOG_PORT"`
	ProbeConfigPath string `mapstructure:"PROBE_CONFIG_PATH"`
}

type ProbeConfig struct {
	Probes []Probe `yaml:"probes"`
}

type Probe struct {
	Name string `yaml:"name" validate:"required"`
	Syslog *SyslogConfig `yaml:"syslog"`
}

type ProbeHandler interface {
	ServeHTTP(*gin.Context)
}

type SyslogConfig struct {
	Address string `yaml:"address" validate:"required"`
	Port int `yaml:"port"`
	// Valid values: rfc5424, rfc3164
	Format string `yaml:"format" validate:"oneof=rfc5424 rfc3164"`
}

func getOutboundIP() string {
	ifaddrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, ifaddr := range ifaddrs {
		netIP, ok := ifaddr.(*net.IPNet)
		if ok && !netIP.IP.IsLoopback() && netIP.IP.To4() != nil {
			return netIP.IP.String()
		}
	}
	return ""
}

func initConfig() {
	viper.BindEnv("SYSLOG_ADDRESS")
	viper.SetDefault("SYSLOG_PORT", 1514)
	viper.SetDefault("CALLBACK_ADDRESS", getOutboundIP())
	viper.SetDefault("CALLBACK_PORT", 8080)
	viper.SetDefault("PROBE_CONFIG_PATH", "/etc/snooze/probes.yaml")

	viper.AutomaticEnv()
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(config.ProbeConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	var probeConfig ProbeConfig
	if err := yaml.Unmarshal(data, &probeConfig); err != nil {
		log.Fatal(err)
	}
	validate := validator.New()
	if err := validate.Struct(probeConfig); err != nil {
		log.Fatal(err)
	}
	for _, probe := range probeConfig.Probes {
		switch {
		case probe.Syslog != nil:
			probes[probe.Name] = NewSyslogRelay(probe.Name, probe.Syslog)
		}
	}
}
