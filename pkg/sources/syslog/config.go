package syslog

import (
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	InstanceName string `mapstructure:"SYSLOG_INSTANCE_NAME"`
	ListenAddress string `mapstructure:"SYSLOG_LISTEN_ADDRESS"`
	ListenPort int `mapstructure:"SYSLOG_LISTEN_PORT"`
}

var config *Config

func initConfig() {
	viper.SetDefault("SYSLOG_LISTEN_ADDRESS", "0.0.0.0")
	viper.SetDefault("SYSLOG_LISTEN_PORT", 1514)
	viper.SetDefault("SYSLOG_INSTANCE_NAME", "")

	viper.AutomaticEnv()
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
}
