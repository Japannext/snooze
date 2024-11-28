package auth

import (
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var parser = jwt.NewParser()

func getISS(tokenString string) (string, error) {
	claims := jwt.MapClaims{}
	token, _, err := parser.ParseUnverified(tokenString, claims)
	if err != nil {
		return "", fmt.Errorf("while parsing token header: %s", err)
	}
	iss, err := token.Claims.GetIssuer()
	if err != nil {
		return "", fmt.Errorf("no field 'iss' found: %s", err)
	}
	return iss, nil
}

type Claims struct {
	Username string `json:"username"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// Authentication middleware
func Authenticated() func(*gin.Context) {
	return func(c *gin.Context) {

		token := c.GetHeader("X-Snooze-Token")
		if token == "" {
			c.String(http.StatusUnauthorized, "no token X-Snooze-Token found")
			c.Abort()
			return
		}

		iss, err := getISS(token)
		if err != nil {
			c.String(http.StatusBadRequest, "issue finding iss: %s", err)
			c.Abort()
			return
		}

		var claims *Claims
		// OIDC authentication
		if m, ok := oidcByUrl[iss]; ok {
			idToken, err :=  m.VerifyToken(token)
			if err != nil {
				c.String(http.StatusUnauthorized, "error verifying token: %s", err)
				c.Abort()
				return
			}
			if err := idToken.Claims(&claims); err != nil {
				c.String(http.StatusUnauthorized, "error verifying claims: %s", err)
				c.Abort()
				return
			}
		// Local token
		} else if iss == "self" {
			err := engine.Verify(token, claims)
			if err == jwt.ErrSignatureInvalid {
				c.String(http.StatusUnauthorized, "wrong JWT signature: %s", err)
				c.Abort()
				return
			}
			if err != nil {
				c.String(http.StatusBadRequest, "error in JWT token: %s", err)
				c.Abort()
				return
			}
		}

		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

