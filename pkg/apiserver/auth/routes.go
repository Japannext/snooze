package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/japannext/snooze/pkg/apiserver/auth/openidconnect"
	"github.com/japannext/snooze/pkg/apiserver/config"
)

var (
	authConfig *config.AuthConfig
)

func RegisterRoutes(r *gin.Engine) {
	authConfig = config.Auth()

	openidconnect.RegisterRoutes(r)
	r.GET("/api/auth/config", getAuthConfig)
}

func getAuthConfig(c *gin.Context) {
	c.JSON(http.StatusOK, config.Auth())
}
