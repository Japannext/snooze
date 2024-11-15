package mq

import (
    "github.com/spf13/viper"
)

type Config struct {
	URL string `mapstructure:"NATS_URL"`
	Replicas int `mapstructure:"NATS_REPLICAS"`
}

var config Config

func initConfig() {
    v := viper.New()

    // Defaults
    v.BindEnv("NATS_URL")
	v.SetDefault("NATS_REPLICAS", "3")

    v.AutomaticEnv()
    err := v.Unmarshal(&config)
    if err != nil {
        log.Fatal(err)
    }
}
