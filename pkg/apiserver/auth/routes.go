package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/japannext/snooze/pkg/apiserver/auth/openidconnect"
	"github.com/japannext/snooze/pkg/apiserver/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer     trace.Tracer
	authConfig *config.AuthConfig
)

func RegisterRoutes(r *gin.Engine) {
	tracer = otel.Tracer("snooze")
	authConfig = config.Auth()

	openidconnect.RegisterRoutes(r)
	r.GET("/api/auth/config", getAuthConfig)
}

func getAuthConfig(c *gin.Context) {
	c.JSON(http.StatusOK, config.Auth())
}
