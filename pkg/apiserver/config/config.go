package config

import (
	"net/url"

	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/common/utils"
)

type EnvConfig struct {
	// Address the HTTP server should listen to
	ListenAddress string `mapstructure:"APISERVER_LISTEN_ADDRESS"`
	// Port the HTTP server should listen to
	ListenPort int `mapstructure:"APISERVER_LISTEN_PORT"`

	// A URL for Jaeger/tracing app, so long as it reacts to /trace/<traceID>
	TraceURL string `mapstructure:"TRACE_URL"`

	// External URL used to expose the apiserver
	ExternalURL url.URL `mapstructure:"APISERVER_EXTERNAL_URL"`

	// The domain name to use for the cookies. Default to the externalURL host.
	CookieDomain string `mapstructure:"APISERVER_COOKIE_DOMAIN"`

	// Authentication config path
	AuthConfig string `mapstructure:"APISERVER_AUTH_CONFIG"`

	// A secret key used to encrypt/decrypt JWT tokens used for authentication
	SecretKey string `mapstructure:"APISERVER_SECRET_KEY"`
}

var env *EnvConfig

func Init() {
	// Defaults
	viper.SetDefault("APISERVER_LISTEN_ADDRESS", "0.0.0.0")
	viper.SetDefault("APISERVER_LISTEN_PORT", 8080)
	viper.SetDefault("APISERVER_AUTH_CONFIG", "/etc/snooze-apiserver/auth.yaml")
	viper.BindEnv("APISERVER_SECRET_KEY")
	viper.AutomaticEnv()
	hooks := []viper.DecoderConfigOption{
		utils.URLDecodeHook(),
	}
	if err := viper.Unmarshal(&env, hooks...); err != nil {
		log.Fatal(err)
	}

	if env.CookieDomain == "" {
		env.CookieDomain = env.ExternalURL.Hostname()
	}

	initAuthConfig()
}

func Env() *EnvConfig {
	if env == nil {
		log.Fatalf("config not initialized!")
	}
	return env
}
