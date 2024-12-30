package openidconnect

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/apiserver/config"
)

var (
	oidcConfig *config.OidcConfig
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	oauth2Config *oauth2.Config
	cookieDomain string

	usernameField string
	emailField string
	fullnameField string
	rolesField string
)

func initConfig() {
	oidcConfig = config.Auth().Oidc
	cookieDomain = config.Env().CookieDomain

	var err error
	provider, err = oidc.NewProvider(context.Background(), oidcConfig.URL)
	if err != nil {
		log.Fatalf("failed to load oidc provider: %s", err)
	}

	oauth2Config = &oauth2.Config{
		ClientID: oidcConfig.ClientID,
		ClientSecret: oidcConfig.ClientSecret,
		RedirectURL: oidcConfig.CallbackURL,
		Endpoint: provider.Endpoint(),
		Scopes: oidcConfig.Scopes,
	}
	verifier = provider.Verifier(&oidc.Config{ClientID: oidcConfig.ClientID})

	usernameField = oidcConfig.UsernameField
	emailField = oidcConfig.EmailField
	fullnameField = oidcConfig.FullnameField
	rolesField = oidcConfig.RolesField
}
