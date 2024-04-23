package otel

import (
  log "github.com/sirupsen/logrus"
  "github.com/spf13/viper"
)

type Config struct {
  // IP address the gRPC server should bind to for the otel listener
  GrpcListeningAddress string `mapstructure:"OTEL_GRPC_LISTENING_ADDRESS"`
  // Port number the gRPC server should bind to for the otel listener
  GrpcListeningPort int `mapstructure:"OTEL_GRPC_LISTENING_PORT"`

  // Whether to enable prometheus metrics
  PrometheusEnable bool `mapstructure:"OTEL_PROMETHEUS_ENABLE"`
  // Port the prometheus exporter should listen to
  PrometheusPort int `mapstructure:"OTEL_PROMETHEUS_PORT"`
}

var config *Config

func initConfig() {
  // Defaults
  viper.SetDefault("OTEL_GRPC_LISTENING_ADDRESS", "0.0.0.0")
  viper.SetDefault("OTEL_GRPC_LISTENING_PORT", 4317)
  viper.SetDefault("OTEL_PROMETHEUS_ENABLE", true)
  viper.SetDefault("OTEL_PROMETHEUS_PORT", 9317)

  viper.AutomaticEnv()
  err := viper.Unmarshal(&config)
  if err != nil {
    log.Fatal(err)
  }
}

