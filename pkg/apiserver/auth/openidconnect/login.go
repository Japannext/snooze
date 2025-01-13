package openidconnect

import (
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/japannext/snooze/pkg/common/utils"
	"golang.org/x/oauth2"
)

const (
	stateByteLength         = 96
	nonceByteLength         = 96
	verifierByteLength      = 96
	cookieExpirationSeconds = 30
)

func login(c *gin.Context) {
	var (
		state    = utils.RandomURLSafeBase64(stateByteLength)
		nonce    = utils.RandomURLSafeBase64(nonceByteLength)
		verifier = utils.RandomURLSafeBase64(verifierByteLength)
	)

	c.SetCookie("state", state, cookieExpirationSeconds, "/", cookieDomain, true, true)
	c.SetCookie("nonce", nonce, cookieExpirationSeconds, "/", cookieDomain, true, true)
	c.SetCookie("verifier", verifier, cookieExpirationSeconds, "/", cookieDomain, true, true)

	url := oauth2Config.AuthCodeURL(state, oidc.Nonce(nonce), oauth2.S256ChallengeOption(verifier))
	c.Redirect(http.StatusTemporaryRedirect, url)
}
