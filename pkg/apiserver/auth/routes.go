package auth

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	initConfig()

	// r.POST("/api/auth/login", postLogin)

	// OIDC routes
	r.GET("/api/oidc/redirect", getOIDCRedirect)
	r.GET("/api/oidc/callback", getOIDCCallback)

	r.GET("/api/auth/methods", getAuthMethods)
}

func getOIDCRedirect(c *gin.Context) {
	method := c.Query("method")
	if method == "" {
		c.String(http.StatusBadRequest, "parameters `method` is required")
		return
	}
	backend, ok := oidcMethods[method]
	if !ok {
		c.String(http.StatusNotFound, "unknown authentication method '%s'", method)
		return
	}

	backend.Redirect(c)
}

func getOIDCCallback(c *gin.Context) {
	method := c.Query("method")
	if method == "" {
		c.String(http.StatusBadRequest, "parameters `method` is required")
		return
	}
	backend, ok := oidcMethods[method]
	if !ok {
		c.String(http.StatusNotFound, "unknown authentication method '%s'", method)
		return
	}

	backend.Callback(c)
}

func getAuthMethods(c *gin.Context) {
	c.JSON(http.StatusOK, authConfig.Methods)
}
