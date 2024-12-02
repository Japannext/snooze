package auth

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"

	"github.com/japannext/snooze/pkg/apiserver/sessions"
)

// Return a random state string.
// It will be a random 64-bit bytes encoded
// in base64 (URL safe)
func getState(c *gin.Context) (string, error) {
	state := c.Query("state")
	if state != "" {
		return state, nil
	}
	nonceBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, nonceBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(nonceBytes), nil
}

func getAuthProviderLogin(c *gin.Context) {
	providerName := c.Param("provider")
	provider, err := goth.GetProvider(providerName)
	if err != nil {
		c.String(http.StatusNotFound, "error getting provider '%s': %s", providerName, err)
		return
	}
	state, err := getState(c)
	if err != nil {
		c.String(http.StatusNotFound, "error getting state: %s", err)
		return
	}
	sess, err := provider.BeginAuth(state)
	if err != nil {
		c.String(http.StatusUnauthorized, "error authenticating: %s", err)
		return
	}
	url, err := sess.GetAuthURL()
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to get authentication URL: %s", err)
		return
	}

	session := sessions.MySession(c)
	session.AuthProvider = providerName
	session.AuthSession = []byte(sess.Marshal())

	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Find the original state based on the auth URL of the goth session
func getAuthURLState(sess goth.Session) (string, error) {
	rawAuthURL, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}
	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		return "", err
	}
	state := authURL.Query().Get("state")
	return state, nil
}

func getAuthProviderCallback(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "getAuthProviderCallback")
	defer span.End()

	providerName := c.Param("provider")
	provider, err := goth.GetProvider(providerName)
	if err != nil {
		c.String(http.StatusNotFound, "error getting provider '%s': %s", providerName, err)
		return
	}

	// Get session
	session := sessions.MySession(c)
	sess, err := provider.UnmarshalSession(string(session.AuthSession))
	if err != nil {
		c.String(http.StatusInternalServerError, "error getting session: %s", err)
		return
	}

	// If user is already authenticated, return it.
	user, err := provider.FetchUser(sess)
	if err == nil {
		c.JSON(http.StatusOK, user)
		return
	}

	// Validate state
	reqState := c.Param("state")
	originalState, err := getAuthURLState(sess)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to retrieve state from auth URL: %s", err)
		return
	}
	if originalState != "" && (originalState != reqState) {
		c.String(http.StatusUnauthorized, "state token mismatch")
		return
	}

	_, err = sess.Authorize(provider, c.Request.URL.Query())
	if err != nil {
		c.String(http.StatusUnauthorized, "error authenticating: %s", err)
		return
	}

	session.Authenticated = true
	if err := session.Save(); err !=  nil {
		c.String(http.StatusInternalServerError, "error saving session: %s", err)
		return
	}

	gu, err := provider.FetchUser(sess)
	if err != nil {
		c.String(http.StatusInternalServerError, "error fetching user info: %s", err)
		return
	}

	c.JSON(http.StatusOK, gu)
}

func getAuthProviderLogout(c *gin.Context) {
	session := sessions.MySession(c)
	session.Delete()
	gothic.Logout(c.Writer, c.Request)
}
