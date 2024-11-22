package auth

import (
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var parser = jwt.NewParser()

func getISS(tokenString string) (string, error) {
	claims := Claims{}
	token, _, err := parser.ParseUnverified(tokenString, claims)
	if err != nil {
		return "", err
	}
	issValue, ok := token.Header["iss"]
	if !ok {
		return "", fmt.Errorf("no field 'iss' found")
	}
	iss, ok := issValue.(string)
	if !ok {
		return "", fmt.Errorf("field 'iss' is not a string.")
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

		token, err := c.Cookie("token")
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		iss, err := getISS(token)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		var claims *Claims
		// OIDC authentication
		if m, ok := oidcMethods[iss]; ok {
			idToken, err :=  m.VerifyToken(token)
			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, err)
				return
			}
			if err := idToken.Claims(&claims); err != nil {
				c.AbortWithError(http.StatusUnauthorized, err)
				return
			}
		// Local token
		} else if iss == "self" {
			err := engine.Verify(token, claims)
			if err == jwt.ErrSignatureInvalid {
				c.AbortWithError(http.StatusUnauthorized, err)
				return
			}
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
		}

		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

