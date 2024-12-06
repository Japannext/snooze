package auth

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/pkg/apiserver/auth/openidconnect"
	"github.com/japannext/snooze/pkg/apiserver/config"
)

var tracer trace.Tracer

func RegisterRoutes(r *gin.Engine) {
	tracer = otel.Tracer("snooze")

	openidconnect.RegisterRoutes(r)
	r.GET("/api/auth/config", getAuthConfig)
}

func getAuthConfig(c *gin.Context) {
	c.JSON(http.StatusOK, config.Auth())
}
