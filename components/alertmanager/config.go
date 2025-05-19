package alertmanager

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	InstanceName  string `mapstructure:"INSTANCE_NAME"`
	ListenAddress string `mapstructure:"LISTEN_ADDRESS"`
	ListenPort    int    `mapstructure:"LISTEN_PORT"`
}

var config *Config

func initConfig() {
	viper.SetDefault("LISTEN_ADDRESS", "0.0.0.0")
	viper.SetDefault("LISTEN_PORT", 9093)
	viper.SetDefault("INSTANCE_NAME", "")

	viper.AutomaticEnv()
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
}
