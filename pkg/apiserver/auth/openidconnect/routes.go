package openidconnect

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	initConfig()

	r.GET("/api/auth/oidc/login", login)
	r.GET("/api/auth/oidc/callback", callback)
	r.GET("/api/auth/oidc/logout", logout)
}
