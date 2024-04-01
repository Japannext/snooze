package logging

import (
  "github.com/spf13/viper"
)

type Config struct {
  Level string `mapStructure:"LEVEL"`
}

func initConfig() (*Config, error) {
  v := viper.New()
  // Prefix
  v.SetEnvPrefix("LOG")

  // Defaults
  v.SetDefault("LEVEL", "debug")

  v.AutomaticEnv()
  cfg := &Config{}
  err := v.Unmarshal(&cfg)
  return cfg, err
}


