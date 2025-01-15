package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/japannext/snooze/pkg/apiserver/token"
	"go.opentelemetry.io/otel"
)

func Authenticate() func(*gin.Context) {
	tracer := otel.Tracer("snooze")

	return func(c *gin.Context) {
		_, span := tracer.Start(c.Request.Context(), "Authenticate")
		defer span.End()

		snoozeToken, err := c.Cookie("snooze-token")
		if err != nil {
			c.String(http.StatusUnauthorized, "no snooze-token found: %s", err)
			c.Abort()
			return
		}

		claims, err := token.Verify(snoozeToken)
		if err != nil {
			c.String(http.StatusUnauthorized, "could not verify token: %s", err)
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("fullname", claims.FullName)
		c.Set("email", claims.Email)
		c.Set("roles", claims.Roles)
	}
}
