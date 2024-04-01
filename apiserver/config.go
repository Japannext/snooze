package apiserver

import (
  "github.com/spf13/viper"
)

type Config struct {
  // Address the HTTP server should listen to
  ListenAddress string `mapstructure:"LISTEN_ADDRESS"`
  // Port the HTTP server should listen to
  ListenPort int `mapstructure:"LISTEN_PORT"`
  // Whether to enable prometheus metrics
  PrometheusEnable bool `mapstructure:"PROMETHEUS_ENABLE"`
  // Port the prometheus exporter should listen to
  PrometheusPort int `mapstructure:"PROMETHEUS_PORT"`
  // Path to the static web pages to serve the snooze webUI
  StaticPath string `mapstructure:"STATIC_PATH"`
}

func (cfg *Config) init() error {
  // Prefix
  viper.SetEnvPrefix("SNOOZE_API")

  // Defaults
  viper.SetDefault("LISTEN_ADDRESS", "0.0.0.0")
  viper.SetDefault("LISTEN_PORT", 8080)
  viper.SetDefault("PROMETHEUS_ENABLE", true)
  viper.SetDefault("PROMETHEUS_PORT", 9080)
  viper.SetDefault("STATIC_PATH", "/static")

  viper.AutomaticEnv()
  return viper.Unmarshal(&config)
}

var config *Config
