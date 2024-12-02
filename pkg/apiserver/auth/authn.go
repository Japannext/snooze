package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/pkg/apiserver/sessions"
)

func Authenticate() func(*gin.Context) {
	return func(c *gin.Context) {
		session := sessions.MySession(c)
		if session.Authenticated {
			c.Next()
			return
		}

		c.String(http.StatusUnauthorized, "session not authenticated")
		c.Abort()
	}
}
