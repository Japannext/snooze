package activecheck

import (
	"net"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"github.com/go-playground/validator/v10"
)

var config *Config
var checks = map[string]Check{}

type Config struct {
	CallbackAddress string `mapstructure:"CALLBACK_ADDRESS"`
	CallbackPort int `mapstructure:"CALLBACK_PORT"`
	SyslogAddress string `mapstructure:"SYSLOG_ADDRESS"`
	SyslogPort int `mapstructure:"SYSLOG_PORT"`
	CheckConfigPath string `mapstructure:"CHECK_CONFIG_PATH"`
}

type CheckConfig struct {
	Checks []Check `yaml:"checks"`
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

	data, err := os.ReadFile(config.CheckConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	var checkConfig CheckConfig
	if err := yaml.Unmarshal(data, &checkConfig); err != nil {
		log.Fatal(err)
	}
	validate := validator.New()
	if err := validate.Struct(checkConfig); err != nil {
		log.Fatal(err)
	}
	for _, check := range checkConfig.Checks {
		check.Load()
		checks[check.Name] = check
	}
}
