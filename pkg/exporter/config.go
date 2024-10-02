package exporter

import (
	"net"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"github.com/go-playground/validator/v10"
)

var config *Config
var checks *CheckConfig

type Config struct {
	CallbackAddress string `mapstructure:"CALLBACK_ADDRESS"`
	CallbackPort int `mapstructure:"CALLBACK_PORT"`
	SyslogAddress string `mapstructure:"SYSLOG_ADDRESS"`
	SyslogPort int `mapstructure:"SYSLOG_PORT"`
	ConfigPath string `mapstructure:"CONFIG_PATH"`
}

type CheckConfig struct {
	SyslogRelays []*SyslogRelay `yaml:"syslog_relays"`
}

func (cc *CheckConfig) FillDefaults() {
	for _, relay := range cc.SyslogRelays {
		relay.FillDefaults()
	}
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
	viper.SetDefault("CONFIG_PATH", "/etc/snooze/activecheck.yaml")

	viper.AutomaticEnv()
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(config.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := yaml.Unmarshal(data, &checks); err != nil {
		log.Fatal(err)
	}
	checks.FillDefaults()

	validate := validator.New()
	if err := validate.Struct(checks); err != nil {
		log.Fatal(err)
	}
}
