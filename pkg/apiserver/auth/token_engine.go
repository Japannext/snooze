package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const X_SNOOZE_TOKEN = "X-Snooze-Token"

type TokenEngine struct {
	secret string
	expiration time.Duration
}

func NewTokenEngine(secret string) *TokenEngine {
	return &TokenEngine{secret: secret}
}

func (engine *TokenEngine) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(engine.secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (engine *TokenEngine) keyfunc(token *jwt.Token) (any, error) {
	return []byte(engine.secret), nil
}

// Verify the token and extract the claims object
func (engine *TokenEngine) Verify(tokenString string, claims jwt.Claims) error {
	opts := []jwt.ParserOption{
		jwt.WithExpirationRequired(),
		jwt.WithIssuer("self"),
	}
	_, err := jwt.ParseWithClaims(tokenString, claims, engine.keyfunc, opts...)
	if err != nil {
		return err
	}
	return nil
}

var engine *TokenEngine

func initTokenEngine(secretKey string) {
	engine = &TokenEngine{
		secret: secretKey,
		expiration: 1 * time.Hour,
	}
}
