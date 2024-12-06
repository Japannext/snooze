package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/japannext/snooze/pkg/common/utils"
)

type AuthConfig struct {
	Oidc *OidcConfig `yaml:"oidc" json:"oidc"`
}

type OidcConfig struct {
	URL string `yaml:"url" json:"url"`
	ClientID string `yaml:"client_id" json:"clientID"`
	ClientSecret string `yaml:"client_secret" json:"-"`
	CallbackURL string `yaml:"callback_url" json:"callbackURL"`
	Scopes []string `yaml:"scopes" json:"scopes"`
	TLSConfig *utils.TLSConfig

	// Cosmetics
	DisplayName string `yaml:"display_name" json:"displayName"`
	Icon string `yaml:"icon" json:"icon"`
	Color string `yaml:"color" json:"color"`
}

type GithubConfig struct {
	ClientID string `yaml:"client_id" json:"clientID"`
	ClientSecret string `yaml:"client_secret" json:"-"`
	CallbackURL string `yaml:"callback_url" json:"callbackURL"`
}

var authConfig *AuthConfig

func Auth() *AuthConfig {
	if authConfig == nil {
		log.Fatalf("auth config is not initialized!")
	}
	return authConfig
}

func initAuthConfig() {
	data, err := os.ReadFile(env.AuthConfig)
	if err != nil {
		log.Fatalf("filed to read config file '%s': %s", env.AuthConfig, err)
	}
	if err := yaml.Unmarshal(data, &authConfig); err != nil {
		log.Fatalf("failed to load config: %s", err)
	}
}
