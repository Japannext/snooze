package config

import (
	"os"

	"github.com/japannext/snooze/pkg/common/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type AuthConfig struct {
	Oidc      *OidcConfig `json:"oidc"      yaml:"oidc"`
	AdminRole string      `json:"adminRole" yaml:"admin_role"`
	UserRole  string      `json:"userRole"  yaml:"user_role"`
}

type OidcConfig struct {
	URL           string   `json:"url"           yaml:"url"`
	ClientID      string   `json:"clientID"      yaml:"client_id"`
	ClientSecret  string   `json:"-"             yaml:"client_secret"`
	CallbackURL   string   `json:"callbackURL"   yaml:"callback_url"`
	Scopes        []string `json:"scopes"        yaml:"scopes"`
	UsernameField string   `json:"usernameField" yaml:"username_field"`
	EmailField    string   `json:"emailField"    yaml:"email_field"`
	FullnameField string   `json:"fullnameField" yaml:"fullname_field"`
	RolesField    string   `json:"rolesField"    yaml:"roles_field"`
	TLSConfig     *utils.TLSConfig

	// Cosmetics
	DisplayName string `json:"displayName" yaml:"display_name"`
	Icon        string `json:"icon"        yaml:"icon"`
	Color       string `json:"color"       yaml:"color"`
}

type GithubConfig struct {
	ClientID     string `json:"clientID"    yaml:"client_id"`
	ClientSecret string `json:"-"           yaml:"client_secret"`
	CallbackURL  string `json:"callbackURL" yaml:"callback_url"`
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

	// Defaults
	if authConfig.Oidc != nil {
		if authConfig.Oidc.UsernameField == "" {
			authConfig.Oidc.UsernameField = "preferred_username"
		}
		if authConfig.Oidc.EmailField == "" {
			authConfig.Oidc.EmailField = "email"
		}
		if authConfig.Oidc.FullnameField == "" {
			authConfig.Oidc.FullnameField = "name"
		}
		if authConfig.Oidc.RolesField == "" {
			authConfig.Oidc.RolesField = "groups"
		}
	}
	if authConfig.AdminRole == "" {
		authConfig.AdminRole = "admin"
	}
	if authConfig.UserRole == "" {
		authConfig.UserRole = "user"
	}
}
