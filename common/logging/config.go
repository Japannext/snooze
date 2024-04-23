package logging

import (
	"github.com/spf13/viper"
)

type Config struct {
	Level string `mapStructure:"LOG_LEVEL"`
}

func initConfig() (*Config, error) {
	v := viper.New()

	// Defaults
	v.SetDefault("LOG_LEVEL", "debug")

	v.AutomaticEnv()
	cfg := &Config{}
	err := v.Unmarshal(&cfg)
	return cfg, err
}
