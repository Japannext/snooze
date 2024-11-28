package auth

import (
	"os"
	"net/url"

	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	// External URL used to expose the apiserver
	ExternalURL string `mapstructure:"APISERVER_EXTERNAL_URL"`
	AuthConfig string `mapstructure:"APISERVER_AUTH_CONFIG"`
	// A secret key used to encrypt/decrypt JWT tokens used for authentication
	SecretKey string `mapstructure:"APISERVER_SECRET_KEY"`
}

type AuthConfig struct {
	GenericOidc *OidcConfig `yaml:"generic_oidc"`
	Methods []AuthMethod `yaml:"methods"`
}

type AuthMethod struct {
    Name string `yaml:"name" json:"name"`
	Kind string `yaml:"kind" json:"kind"`
    DisplayName string `yaml:"display_name" json:"displayName"`
    Icon string `yaml:"icon" json:"icon"`
    Color string `yaml:"color" json:"color"`
	// Methods
    Oidc *OidcMethod `yaml:"oidc" json:"oidc"`
}

var config *Config
var authConfig *AuthConfig
var oidcMethods = make(map[string]*OidcMethod)
var oidcByUrl = make(map[string]*OidcMethod)
var cookieDomain string

func initConfig() {
	v := viper.New()
	v.SetDefault("APISERVER_AUTH_CONFIG", "/etc/snooze-apiserver/auth.yaml")
	v.AutomaticEnv()
	if err := v.Unmarshal(&config); err != nil {
		log.Fatalf("failed to load config: %s", err)
	}
	data, err := os.ReadFile(config.AuthConfig)
	if err != nil {
		log.Fatalf("failed to read config at %s: %s", config.AuthConfig, err)
	}
	if err := yaml.Unmarshal(data, &authConfig); err != nil {
		log.Fatalf("failed to unmarshal config at %s: %s", config.AuthConfig, err)
	}

	u, err := url.Parse(config.ExternalURL)
	if err != nil {
		log.Fatalf("failed to parse URL '%s': %s", config.ExternalURL, err)
	}
	cookieDomain = u.Hostname()

	// To verify the uniqueness of names
	var uniq = map[string]bool{}
	for _, method := range authConfig.Methods {
		if _, ok := uniq[method.Name]; ok {
			log.Fatalf("Duplicate name `%s` found in auth config %s", method.Name, config.AuthConfig)
		}
		switch(method.Kind) {
		case "oidc":
			if method.Oidc == nil {
				log.Fatalf("oidc options not defined for '%s'", method.Name)
			}
			if err := method.Oidc.Load(); err != nil {
				log.Fatalf("failed to load '%s' OIDC backend: %s", method.Name, err)
			}
			oidcMethods[method.Name] = method.Oidc
			oidcByUrl[method.Oidc.URL] = method.Oidc
		}
		log.Infof("loaded auth method '%s'", method.Name)
	}

	initTokenEngine(config.SecretKey)
}
