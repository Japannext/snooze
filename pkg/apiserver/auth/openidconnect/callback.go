package openidconnect

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/japannext/snooze/pkg/apiserver/token"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func callback(c *gin.Context) {
	state, err := c.Cookie("state")
	if err != nil {
		log.Warnf("cookie not found in callback: %s", err)
		c.String(http.StatusUnauthorized, "state cookie not found: %s", err)
		return
	}
	if c.Query("state") != state {
		log.Warnf("state cookie and parameter do not match")
		c.String(http.StatusUnauthorized, "state cookie and parameter do not match")
		return
	}
	if errString := c.Query("error"); errString != "" {
		desc := c.Query("error_description")
		log.Warnf("OIDC callback returned error: %s | %s", errString, desc)
		c.String(http.StatusUnauthorized, "OIDC callback returned error: %s | %s", errString, desc)
		return
	}
	pkceVerifier, err := c.Cookie("verifier")
	if err != nil {
		log.Warnf("verifier not found in callback: %s", err)
		c.String(http.StatusUnauthorized, "verifier cookie not found: %s", err)
		return
	}

	ctx := c.Request.Context()
	// ctx = oidc.ClientContext(ctx, &http.Client{Transport: &http.Transport{TLSClientConfig: m.internal.tlsConfig}})

	oauth2Token, err := oauth2Config.Exchange(ctx, c.Query("code"), oauth2.VerifierOption(pkceVerifier))
	if err != nil {
		c.String(http.StatusUnauthorized, "Failed to exchange token: %s", err)
		return
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		c.String(http.StatusUnauthorized, "No id_token field in oauth2 token: %+v", oauth2Token)
		return
	}
	idToken, err := verifier.Verify(c, rawIDToken)
	if err != nil {
		c.String(http.StatusUnauthorized, "Failed to verify ID token: %s", err)
		return
	}
	nonce, err := c.Cookie("nonce")
	if err != nil {
		c.String(http.StatusUnauthorized, "nonce not found: %s", err)
		return
	}
	if idToken.Nonce != nonce {
		c.String(http.StatusUnauthorized, "nonce did not match")
		return
	}

	// Claims
	claims := make(map[string]interface{})
	if err := idToken.Claims(&claims); err != nil {
		c.String(http.StatusInternalServerError, "error verifying claims: %s", err)
		return
	}
	delete(claims, "nonce")

	snoozeAuthClaims := token.NewClaims()
	snoozeAuthClaims.Username = extractString(claims, usernameField)
	snoozeAuthClaims.Email = extractString(claims, emailField)
	snoozeAuthClaims.FullName = extractString(claims, fullnameField)
	snoozeAuthClaims.Roles = extractStringSlice(claims, rolesField)

	snoozeToken, err := token.Sign(snoozeAuthClaims)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to generate snooze-token: %s", err)
		return
	}
	c.SetCookie("snooze-token", snoozeToken, 3600, "/", cookieDomain, true, true)
}

func extractString(claims map[string]interface{}, field string) string {
	data, found := claims[field]
	if !found {
		log.Debugf("field '%s' not found in claim", field)
		return ""
	}
	str, ok := data.(string)
	if !ok {
		log.Debugf("field '%s' is not a string", field)
		return ""
	}
	return str
}

func extractStringSlice(claims map[string]interface{}, field string) []string {
	data, found := claims[field]
	if !found {
		log.Debugf("field '%s' not found in claim", field)
		return []string{}
	}
	slice, ok := data.([]interface{})
	if !ok {
		log.Debugf("field '%s' is not a slice: %T", field, data)
		return []string{}
	}
	var results []string
	for _, e := range slice {
		str, ok := e.(string)
		if !ok {
			log.Debugf("field '%s': ignoring non-string value (%T): %+v", field, e, e)
			continue
		}
		results = append(results, str)
	}
	return results
}
