package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, role := range c.GetStringSlice("roles") {
			if role == authConfig.AdminRole {
				c.Next()
				return
			}
		}
		c.String(http.StatusForbidden, "You are not allowed on this page. Need to be admin")
		c.Abort()
		return
	}
}

func UserAllowed() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, role := range c.GetStringSlice("roles") {
			if role == authConfig.AdminRole || role == authConfig.UserRole {
				c.Next()
				return
			}
		}
		c.String(http.StatusForbidden, "You are not allowed on this page. Need to be user or admin")
		c.Abort()
		return
	}
}
