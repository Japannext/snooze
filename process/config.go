package server

import (
  log "github.com/sirupsen/logrus"
  "github.com/spf13/viper"
)

type ConfigModel struct {
  // IP address the gRPC server should bind to for the otel listener
  OtelGrpcListeningAddress string `mapstructure:"OTEL_GRPC_LISTENING_ADDRESS"`
  // Port number the gRPC server should bind to for the otel listener
  OtelGrpcListeningPort int `mapstructure:"GRPC_LISTENING_PORT"`
  // A logrus log level (trace/debug/info/warning/error/fatal/panic).
  LogLevel string `mapstructure:"LOG_LEVEL"`
  // Whether to enable prometheus metrics
  PrometheusEnable bool `mapstructure:"PROMETHEUS_ENABLE"`
  // Port the prometheus exporter should listen to
  PrometheusPort int `mapstructure:"PROMETHEUS_PORT"`
}

var Config ConfigModel

func setDefaults() {
  viper.SetDefault("OTEL_GRPC_LISTENING_ADDRESS", "0.0.0.0")
  viper.SetDefault("OTEL_GRPC_LISTENING_PORT", 4317)
  viper.SetDefault("LOG_LEVEL", "debug")
  viper.SetDefault("PROMETHEUS_ENABLE", true)
  viper.SetDefault("PROMETHEUS_PORT", 9317)
}

func initConfig() error {
  viper.SetEnvPrefix("SNOOZE")
  setDefaults()
  viper.AutomaticEnv()
  if err = viper.Unmarshal(&Config); err != nil {
    return err
  }
  log.Debugf("Loaded config: %+v", Config)
  return nil
}
