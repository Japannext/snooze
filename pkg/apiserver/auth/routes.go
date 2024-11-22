package auth

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	r.POST("/api/oidc/login/:method", postLogin)
	r.POST("/api/oidc/callback/:method", postLoginCallback)

	r.GET("/api/auth/methods", getAuthMethods)
}

func postLogin(c *gin.Context) {
	method := c.Param("method")
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

func postLoginCallback(c *gin.Context) {
	method := c.Param("method")
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
