package rabbitmq

import (
	"github.com/spf13/viper"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Address string `mapstructure:"AMQP_ADDRESS"`
	Username string `mapstructure:"AMQP_USERNAME"`
	Password string `mapstructure:"AMQP_PASSWORD"`
}

func initConfig() (string, *amqp.Config) {
	v := viper.New()

	// Defaults
	v.SetDefault("AMQP_ADDRESS", "amqp://127.0.0.1:5672/")
	v.BindEnv("AMQP_USERNAME")
	v.BindEnv("AMQP_PASSWORD")

	v.AutomaticEnv()
	var cfg *Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Creating AMQP config
	auth := &amqp.PlainAuth{cfg.Username, cfg.Password}
	config := &amqp.Config{
		SASL: []amqp.Authentication{auth},
	}
	return cfg.Address, config
}
