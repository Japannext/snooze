package auth

import (
	"os"
	"net/url"

	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/openidConnect"

	"github.com/japannext/snooze/pkg/common/utils"
)

var config *Config
var authConfig *AuthConfig
var cookieDomain string

type Config struct {
	// External URL used to expose the apiserver
	ExternalURL string `mapstructure:"APISERVER_EXTERNAL_URL"`
	AuthConfig string `mapstructure:"APISERVER_AUTH_CONFIG"`
	// A secret key used to encrypt/decrypt JWT tokens used for authentication
	SecretKey string `mapstructure:"APISERVER_SECRET_KEY"`
}

type AuthConfig struct {
	Oidc *OidcConfig `yaml:"oidc" json:"oidc,omitempty"`
	Github *GithubConfig `yaml:"github" json:"github,omitempty"`
}

type OidcConfig struct {
	URL string `yaml:"url" json:"url"`
	ClientID string `yaml:"client_id" json:"clientID"`
	ClientSecret string `yaml:"client_secret" json:"-"`
	RedirectURL string `yaml:"redirect_url" json:"redirectURL"`
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

	// Goth
	if cfg := authConfig.Oidc; cfg != nil {
		provider, err := openidConnect.New(cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL, cfg.URL)
		if err != nil {
			log.Fatalf("error with oidc provider: %s", err)
		}
		goth.UseProviders(provider)
	}
	if cfg := authConfig.Github; cfg != nil {
		goth.UseProviders(github.New(cfg.ClientID, cfg.ClientSecret, cfg.CallbackURL, "openid", "profile", "email"))
	}
}
