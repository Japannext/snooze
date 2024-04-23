package rabbitmq

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Address string `mapstructure:"AMQP_ADDRESS"`
}

func initConfig() *Config {
	v := viper.New()

	// Defaults
	v.SetDefault("AMQP_ADDRESS", "amqp://127.0.0.1:5672/")

	v.AutomaticEnv()
	var cfg *Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
