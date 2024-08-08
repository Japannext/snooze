package authentication

import (
	"context"
	"encoding/base64"
	"time"

	"golang.org/x/oauth2"
	"github.com/coreos/go-oidc/v3/oidc"
)

type OidcConfig struct {
	Name string
	URL string
	ClientID string
	ClientSecret string
	RedirectURL string
	Scopes []string
}

type OidcBackend struct {
	Name string
	Provider oidc.Provider
	Verifier oidc.IDTokenVerifier
	Config oauth2.Config
}

func NewOidcBackend(ctx context.Context, cfg OidcConfig) (*OidcBackend, error) {
	provider, err := oidc.NewProvider(ctx, cfg.URL)
	if err != nil {
	}
	config := oauth2.Config{
		ClientID: cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL: cfg.RedirectURL,
		Endpoint: provider.Endpoint(),
		Scopes: cfg.Scopes,
	}
	verifier := provider.Verifier(config)
	return &OidcBackend{cfg.Name, provider, verifier, config}
}

// Generate a random 16 bit string in base64 format
func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func setCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}

func (o *OidcBackend) Redirect(w http.ResponseWriter, r *http.Request) error {
	state, err := randString(16)
	if err != nil {
		log.Error(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	none, err := randString(16)
	if err != nil {
		log.Error(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	setCallbackCookie(w, r, "state", state)
	setCallbackCookie(w, r, "nonce", nonce)

	url, err := o.Config.AuthCodeURL(state, oidc.Nonce(nonce))
	if err != nil {
		log.Error(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (o *OidcBackend) Callback(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	state, err := r.Cookie("state")
	if err != nil {
		http.Error(w, "state not found", http.StatusBadRequest)
		return
	}
	if r.URL.Query().Get("state") != state.Value {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}
	errorCode := r.URL.Query().Get("error")
	if errorCode != "" {
		errorDescription := r.URL.Query().Get("error_description")
		msg := fmt.Sprintf("OIDC callback returned error: %s - %s", errorCode, errorDescription)
		log.Errorf(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	oauth2Token, err := o.Config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token: " + err.Error(), http.StatusInternalServerError)
		return
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}
	idToken, err := o.Verifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: " + err.Error(), http.StatusInternalServerError)
		return
	}

	nonce, err := r.Cookie("nonce")
	if err != nil {
		http.Error(w, "nonce not found", http.StatusBadRequest)
		return
	}
	if idToken.Nonce != nonce.Value {
		http.Error(w, "nonce did not match", http.StatusBadRequest)
		return
	}

	oauth2Token.AccessToken = "*REDACTED*"

	resp := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	}{oauth2Token, new(json.RawMessage)}

	if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (o *OidcBackend) RegisterRoutes(r *gin.Engine) {
	r.HandleFunc(fmt.Sprintf("/oidc/%s/login", o.Name), o.Redirect)
	r.HandleFunc(fmt.Sprintf("/oidc/%s/callback", o.Name), o.Callback)
}
