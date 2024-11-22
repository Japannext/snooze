package auth

import (
	"os"

	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AuthConfig string `mapstructure:"APISERVER_AUTH_CONFIG"`
}

type AuthConfig struct {
	Methods []AuthMethod `yaml:"methods"`
}

type AuthMethod struct {
    Name string `yaml:"name" json:"name"`
    DisplayName string `yaml:"display_name" json:"displayName"`
    Icon string `yaml:"icon" json:"icon"`
    Color string `yaml:"color" json:"color"`
	// Methods
    Oidc *OidcMethod `yaml:"oidc" json:"-"`
}

var config *Config
var authConfig *AuthConfig
var oidcMethods map[string]*OidcMethod

func initConfig() {
	v := viper.New()
	v.SetDefault("APISERVER_AUTH_CONFIG", "/etc/snooze-apiserver/auth_config.yaml")
	v.AutomaticEnv()
	if err := v.Unmarshal(&config); err != nil {
		log.Fatalf("failed to load config: %s", err)
	}
	data, err := os.ReadFile(config.AuthConfig)
	if err != nil {
		log.Fatalf("failed to read config at %s: %s", config.AuthConfig, err)
	}
	if err := yaml.Unmarshal(data, &authConfig); err != nil {
		log.Fatal("failed to unmarshal config at %s: %s", config.AuthConfig, err)
	}
	// To verify the uniqueness of names
	var uniq = map[string]bool{}
	for _, method := range authConfig.Methods {
		if _, ok := uniq[method.Name]; ok {
			log.Fatalf("Duplicate name `%s` found in auth config %s", method.Name, config.AuthConfig)
		}
		if method.Oidc != nil {
			if err := method.Oidc.Load(); err != nil {
				log.Fatalf("failed to load '%s' OIDC backend: %s", method.Name, err)
			}
			oidcMethods[method.Oidc.URL] = method.Oidc
		}
	}
}
