package auth

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"github.com/gin-gonic/gin"
	"github.com/coreos/go-oidc/v3/oidc"
	log "github.com/sirupsen/logrus"
)

type OidcConfig struct {
	ClientID string `yaml:"client_id" json:"clientID"`
	ClientSecret string `yaml:"client_secret" json:"-"`
	RedirectURL string `yaml:"redirect_url" json:"redirectURL"`
	Scopes []string `yaml:"scopes" json:"scopes"`
	CaCert string `yaml:"cacert"`

	// Cosmetics
	DisplayName string `yaml:"name" json:"name"`
	Icon string `yaml:"icon" json:"icon"`
	Color string `yaml:"color" json:"color"`

	internal struct {
		provider *oidc.Provider
		verifier *oidc.IDTokenVerifier
		config oauth2.Config
		tlsConfig *tls.Config
	}
}

type OidcMethod struct {
	URL string `yaml:"url" json:"url"`
	ClientID string `yaml:"client_id" json:"clientID"`
	ClientSecret string `yaml:"client_secret" json:"-"`
	RedirectURL string `yaml:"redirect_url" json:"redirectURL"`
	Scopes []string `yaml:"scopes" json:"scopes"`
	CaCert string `yaml:"cacert"`

	internal struct {
		provider *oidc.Provider
		verifier *oidc.IDTokenVerifier
		config oauth2.Config
		tlsConfig *tls.Config
	}
}

func (m *OidcMethod) Load() error {
	var err error
	ctx := context.Background()
	if m.CaCert != "" {
		caCertData, err := ioutil.ReadFile(m.CaCert)
		if err != nil {
			return err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCertData)
		tlsConfig := &tls.Config{RootCAs: caCertPool}
		m.internal.tlsConfig = tlsConfig
		ctx = oidc.ClientContext(ctx, &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}})
	}
	m.internal.provider, err = oidc.NewProvider(ctx, m.URL)
	if err != nil {
		return err
	}
	m.internal.config = oauth2.Config{
		ClientID: m.ClientID,
		ClientSecret: m.ClientSecret,
		RedirectURL: m.RedirectURL,
		Endpoint: m.internal.provider.Endpoint(),
		Scopes: []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}
	m.internal.verifier = m.internal.provider.Verifier(&oidc.Config{ClientID: m.ClientID})
	return nil
}

// Generate a random 16 bit string in base64 format
func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (m *OidcMethod) VerifyToken(tokenString string) (*oidc.IDToken, error) {
	ctx := context.Background()
	idToken, err := m.internal.verifier.Verify(ctx, tokenString)
	if err != nil {
		return nil, err
	}
	return idToken, nil
}

func (m *OidcMethod) Redirect(c *gin.Context) {
	state, err := randString(16)
	if err != nil {
		log.Error(err)
		c.String(http.StatusInternalServerError, "Error getting a random string for state: %w", err)
		return
	}
	nonce, err := randString(16)
	if err != nil {
		log.Error(err)
		c.String(http.StatusInternalServerError, "Error getting a random string for nonce: %w", err)
		return
	}
	c.SetCookie("state", state, 3600, "/", cookieDomain, true, true)
	c.SetCookie("nonce", nonce, 3600, "/", cookieDomain, true, true)
	url := m.internal.config.AuthCodeURL(state, oidc.Nonce(nonce))
	c.Redirect(http.StatusFound, url)
}

func (m *OidcMethod) Callback(c *gin.Context) {
	state, err := c.Cookie("state")
	if err != nil {
		log.Error(err)
		c.String(http.StatusBadRequest, "state cookie not found: %s", err)
		return
	}
	if c.Query("state") != state {
		c.String(http.StatusBadRequest, "state cookie and parameter do not match")
		return
	}
	if c.Query("error") != "" {
		c.String(http.StatusBadGateway, "OIDC callback returned error: %s | %s", c.Query("error"), c.Query("error_description"))
		return
	}

	ctx := c.Request.Context()
	ctx = oidc.ClientContext(ctx, &http.Client{Transport: &http.Transport{TLSClientConfig: m.internal.tlsConfig}})
	oauth2Token, err := m.internal.config.Exchange(ctx, c.Query("code"))
	if err != nil {
		c.String(http.StatusBadGateway, "Failed to exchange token: %s", err)
		return
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		c.String(http.StatusBadGateway, "No id_token field in oauth2 token: %+v", oauth2Token)
		return
	}
	idToken, err := m.internal.verifier.Verify(c, rawIDToken)
	if err != nil {
		c.String(http.StatusBadGateway, "Failed to verify ID token: %s", err)
		return
	}
	nonce, err := c.Cookie("nonce")
	if err != nil {
		c.String(http.StatusBadGateway, "nonce not found: %s", err)
		return
	}
	if idToken.Nonce != nonce {
		c.String(http.StatusBadGateway, "nonce did not match")
		return
	}
	oauth2Token.AccessToken = "*REDACTED*"
	resp := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	}{oauth2Token, new(json.RawMessage)}

	if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
		c.String(http.StatusInternalServerError, "error verifying claims: %s", err)
		return
	}
	c.SetCookie("token", rawIDToken, 3600, "/", cookieDomain, true, true)
	c.SetCookie("username", "john.doe", 3600, "/", cookieDomain, false, false)
	c.SetCookie("role", "admin", 3600, "/", cookieDomain, false, false)

	c.JSON(http.StatusOK, resp)
}
