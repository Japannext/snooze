package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, role := range ctx.GetStringSlice("roles") {
			if role == authConfig.AdminRole {
				ctx.Next()
				return
			}
		}

		ctx.String(http.StatusForbidden, "You are not allowed on this page. Need to be admin")
		ctx.Abort()
	}
}

func UserAllowed() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, role := range ctx.GetStringSlice("roles") {
			if role == authConfig.AdminRole || role == authConfig.UserRole {
				ctx.Next()
				return
			}
		}

		ctx.String(http.StatusForbidden, "You are not allowed on this page. Need to be user or admin")
		ctx.Abort()
	}
}
