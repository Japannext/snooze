package apiserver

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"github.com/gin-contrib/cors"
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
	// Path to the authentication config (in yaml)
	AuthConfig string `mapstructure:"APISERVER_AUTH_CONFIG"`
	// Path to a CORS configuration
	CorsConfig string `mapstructure:"APISERVER_CORS_CONFIG"`

	// A secret key used to encrypt/decrypt JWT tokens used for authentication
	SecretKey string `mapstructure:"APISERVER_SECRET_KEY"`

	// A URL for Jaeger/tracing app, so long as it reacts to /trace/<traceID>
	TraceURL string `mapstructure:"TRACE_URL"`
}

type AuthConfig struct {
	Backends map[string]AuthBackend `yaml:"backends"`
}

type AuthBackend struct {
	ID string `yaml:"id"`
	DisplayName string `yaml:"display_name"`
	Icon string `yaml:"icon"`
	Color string `yaml:"color"`
	Oidc *OidcConfig `yaml:"oidc"`
	// Static *StaticConfig `yaml:"static"`
}

type CorsConfig struct {
	AllowAllOrigins bool `yaml:"allow_all_origins"`
	AllowCredentials bool `yaml:"allow_credentials"`
	AllowHeaders []string `yaml:"allow_headers"`
}

func loadCorsConfig() {
	data, err := os.ReadFile(config.CorsConfig)
	if err != nil {
		log.Fatal(err)
	}
	cc := &CorsConfig{}
	if err := yaml.Unmarshal(data, &cc); err != nil {
		log.Fatal(err)
	}

	corsConfig.AllowAllOrigins = cc.AllowAllOrigins
	corsConfig.AllowCredentials = cc.AllowCredentials
	corsConfig.AllowHeaders = cc.AllowHeaders
}

var config *Config
var corsConfig *cors.Config
var authConfig *AuthConfig

func initConfig() {
	// Defaults
	viper.SetDefault("APISERVER_LISTEN_ADDRESS", "0.0.0.0")
	viper.SetDefault("APISERVER_LISTEN_PORT", 8080)
	viper.SetDefault("APISERVER_PROMETHEUS_ENABLE", true)
	viper.SetDefault("APISERVER_PROMETHEUS_PORT", 9080)
	viper.SetDefault("APISERVER_STATIC_PATH", "/static")
	viper.SetDefault("APISERVER_AUTH_CONFIG", "/etc/snooze-apiserver/auth_config.yaml")
	viper.SetDefault("APISERVER_CORS_CONFIG", "")

	viper.BindEnv("APISERVER_SECRET_KEY")

	viper.AutomaticEnv()
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	// Load auth backends
	data, err := os.ReadFile(config.AuthConfig)
	if err != nil {
		log.Fatal(err)
	}
	if err := yaml.Unmarshal(data, &authConfig); err != nil {
		log.Fatal(err)
	}

	// CORS
	cc := cors.DefaultConfig()
	cc.AllowAllOrigins = true
	corsConfig = &cc
	if config.CorsConfig != "" {
		loadCorsConfig()
	}
}
