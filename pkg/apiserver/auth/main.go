package auth

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"github.com/gin-gonic/gin"
)

var tracer trace.Tracer

func RegisterAuthRoutes(r *gin.Engine) {
	initConfig()
	tracer = otel.Tracer("snooze")

	r.GET("/api/auth/:provider/login", getAuthProviderLogin)
	r.GET("/api/auth/:provider/callback", getAuthProviderCallback)
	r.GET("/api/auth/:provider/logout", getAuthProviderLogout)
}
