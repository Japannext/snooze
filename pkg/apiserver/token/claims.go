package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims `json:",inline"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	Email string `json:"email"`
	Roles []string `json:"roles"`
}

func NewClaims() *Claims {
	return &Claims{RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(6 * time.Hour)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer: "snooze-apiserver",
		Subject: "snooze",
		Audience: []string{"snooze"},
	}}
}
