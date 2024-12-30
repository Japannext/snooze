package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/japannext/snooze/pkg/apiserver/config"
)

var secret []byte
var expiration time.Duration = time.Duration(6) * time.Hour

func Init() {
	env := config.Env()
	secret = []byte(env.SecretKey)
}

func Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func keyfunc(token *jwt.Token) (any, error) {
	return []byte(secret), nil
}

// Verify the token and extract the claims object
func Verify(tokenString string, ) (*Claims, error) {
	claims := &Claims{}
	opts := []jwt.ParserOption{
		jwt.WithExpirationRequired(),
	}
	_, err := jwt.ParseWithClaims(tokenString, claims, keyfunc, opts...)
	if err != nil {
		return claims, err
	}
	return claims, nil
}
