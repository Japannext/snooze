package apiserver

import (
  log "github.com/sirupsen/logrus"
  "github.com/spf13/viper"
)

type Config struct {
  // Address the HTTP server should listen to
  ListenAddress string `mapstructure:"APISERVER_LISTEN_ADDRESS"`
  // Port the HTTP server should listen to
  ListenPort int `mapstructure:"APISERVER_LISTEN_PORT"`
  // Whether to enable prometheus metrics
  PrometheusEnable bool `mapstructure:"APISERVER_PROMETHEUS_ENABLE"`
  // Port the prometheus exporter should listen to
  PrometheusPort int `mapstructure:"APISERVER_PROMETHEUS_PORT"`
  // Path to the static web pages to serve the snooze webUI
  StaticPath string `mapstructure:"APISERVER_STATIC_PATH"`
}

var config *Config

func initConfig() {

  // Defaults
  viper.SetDefault("APISERVER_LISTEN_ADDRESS", "0.0.0.0")
  viper.SetDefault("APISERVER_LISTEN_PORT", 8080)
  viper.SetDefault("APISERVER_PROMETHEUS_ENABLE", true)
  viper.SetDefault("APISERVER_PROMETHEUS_PORT", 9080)
  viper.SetDefault("APISERVER_STATIC_PATH", "/static")

  viper.AutomaticEnv()
  if err := viper.Unmarshal(&config); err != nil {
    log.Fatal(err)
  }
  log.Debugf("Loaded config: %+v", config)
}
