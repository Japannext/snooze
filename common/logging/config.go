package logging

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Level string `mapstructure:"LOG_LEVEL"`
}

func initConfig() *Config {
	v := viper.New()

	// Defaults
	v.SetDefault("LOG_LEVEL", "info")

	v.AutomaticEnv()
	cfg := &Config{}
	err := v.Unmarshal(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
