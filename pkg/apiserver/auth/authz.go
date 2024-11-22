package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.GetString("role") {
		case "admin":
			c.Next()
		default:
			c.String(http.StatusForbidden, "You are not allowed on this page. Admin only.")
			c.Abort()
		}
	}
}
