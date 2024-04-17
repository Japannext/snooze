package redis

import (
  "github.com/spf13/viper"
  redisv9 "github.com/redis/go-redis/v9"
)

type Config struct {
  Address string `mapstructure:"REDIS_ADDRESS"`
  Password string `mapstructure:"REDIS_PASSWORD"`
  DB int `mapstructure:"REDIS_DB"`
}

func initConfig() (*redisv9.Options, error) {
  v := viper.New()

  // Defaults
  v.SetDefault("REDIS_ADDRESS", "127.0.0.1:6379")

  v.AutomaticEnv()
  cfg := &Config{}
  if err := v.Unmarshal(&cfg); err != nil {
    return &redisv9.Options{}, err
  }
  options := &redisv9.Options{
    Addr: cfg.Address,
    Password: cfg.Password,
    DB: cfg.DB,
  }
  return options, nil
}
