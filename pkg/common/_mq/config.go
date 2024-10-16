package mq

import (
    "github.com/spf13/viper"
)

type Config struct {
	URL string `mapstructure:"NATS_URL"`
}

var config Config

func initConfig() {
    v := viper.New()

    // Defaults
    v.BindEnv("NATS_URL")

    v.AutomaticEnv()
    err := v.Unmarshal(&config)
    if err != nil {
        log.Fatal(err)
    }
}
