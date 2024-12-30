package openidconnect

import (
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/japannext/snooze/pkg/common/utils"
)

func login(c *gin.Context) {
	state, err := utils.RandomURLSafeBase64(16)
	if err != nil {
			log.Error(err)
			c.String(http.StatusInternalServerError, "Error getting a random string for state: %s", err)
			return
	}
	nonce, err := utils.RandomURLSafeBase64(16)
	if err != nil {
			log.Error(err)
			c.String(http.StatusInternalServerError, "Error getting a random string for nonce: %s", err)
			return
	}
	verifier, err := utils.RandomURLSafeBase64(32)
	if err != nil {
		log.Error(err)
		c.String(http.StatusInternalServerError, "Error getting a random string for pkce: %s", err)
	}

	c.SetCookie("state", state, 30, "/", cookieDomain, true, true)
	c.SetCookie("nonce", nonce, 30, "/", cookieDomain, true, true)
	c.SetCookie("verifier", verifier, 30, "/", cookieDomain, true, true)

	url := oauth2Config.AuthCodeURL(state, oidc.Nonce(nonce), oauth2.S256ChallengeOption(verifier))
	c.Redirect(http.StatusTemporaryRedirect, url)
}
