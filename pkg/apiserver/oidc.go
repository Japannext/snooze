package apiserver

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
	URL string `yaml:"url"`
	ClientID string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURL string `yaml:"redirect_url"`
	Scopes []string `yaml:"scopes"`
	CaCert string `yaml:"cacert"`
}

type OidcBackend struct {
	ID string
	Provider *oidc.Provider
	Verifier *oidc.IDTokenVerifier
	Config oauth2.Config
}

func NewOidcBackend(ctx context.Context, id string, cfg *OidcConfig) (*OidcBackend, error) {
	if cfg.CaCert != "" {
		caCertData, err := ioutil.ReadFile(cfg.CaCert)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCertData)
		tlsConfig := &tls.Config{RootCAs: caCertPool}
		ctx = oidc.ClientContext(ctx, &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}})
	}
	provider, err := oidc.NewProvider(ctx, cfg.URL)
	if err != nil {
		return nil, err
	}
	oauth2Config := oauth2.Config{
		ClientID: cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL: cfg.RedirectURL,
		Endpoint: provider.Endpoint(),
		Scopes: cfg.Scopes,
	}
	oidcConfig := &oidc.Config{ClientID: cfg.ClientID}
	verifier := provider.Verifier(oidcConfig)
	return &OidcBackend{id, provider, verifier, oauth2Config}, nil
}

// Generate a random 16 bit string in base64 format
func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (o *OidcBackend) Redirect(c *gin.Context) {
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
	c.SetCookie("state", state, 3600, "/", "localhost", false, true)
	c.SetCookie("nonce", nonce, 3600, "/", "localhost", false, true)
	url := o.Config.AuthCodeURL(state, oidc.Nonce(nonce))
	c.Redirect(http.StatusFound, url)
}

func (o *OidcBackend) Callback(c *gin.Context) {
	state, err := c.Cookie("state")
	if err != nil {
		log.Error(err)
		c.String(http.StatusBadRequest, "state cookie not found: %w", err)
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
	oauth2Token, err := o.Config.Exchange(c, c.Query("code"))
	if err != nil {
		c.String(http.StatusBadGateway, "Failed to exchange token: %w", err)
		return
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		c.String(http.StatusBadGateway, "No id_token field in oauth2 token: %w", err)
		return
	}
	idToken, err := o.Verifier.Verify(c, rawIDToken)
	if err != nil {
		c.String(http.StatusBadGateway, "Failed to verify ID token: %w", err)
		return
	}
	nonce, err := c.Cookie("nonce")
	if err != nil {
		c.String(http.StatusBadGateway, "nonce not found: %w", err)
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
		c.String(http.StatusInternalServerError, "error verifying claims: %w", err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
