package otel

import (
  "github.com/spf13/viper"
)

type Config struct {
  // IP address the gRPC server should bind to for the otel listener
  GrpcListeningAddress string `mapstructure:"GRPC_LISTENING_ADDRESS"`
  // Port number the gRPC server should bind to for the otel listener
  GrpcListeningPort int `mapstructure:"GRPC_LISTENING_PORT"`

  // Whether to enable prometheus metrics
  PrometheusEnable bool `mapstructure:"PROMETHEUS_ENABLE"`
  // Port the prometheus exporter should listen to
  PrometheusPort int `mapstructure:"PROMETHEUS_PORT"`
}

func (cfg *Config) init() error {
  // Prefix
  viper.SetEnvPrefix("SNOOZE_OTEL")
  // Defaults
  viper.SetDefault("GRPC_LISTENING_ADDRESS", "0.0.0.0")
  viper.SetDefault("GRPC_LISTENING_PORT", 4317)
  viper.SetDefault("PROMETHEUS_ENABLE", true)
  viper.SetDefault("PROMETHEUS_PORT", 9317)

  viper.AutomaticEnv()
  return viper.Unmarshal(&config)
}

var config *Config
