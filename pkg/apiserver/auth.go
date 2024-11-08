package apiserver

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func init() {
	routes = append(routes, registerAuthRoutes)
}

type AuthButton struct {
    DisplayName string `json:"displayName"`
    Icon string `json:"icon"`
	Color string `json:"color"`
	Path string `json:"path"`
}

type AuthButtons struct {
	Buttons []AuthButton `json:"buttons"`
}

var authButtons []AuthButton

func listAuthButtons(c *gin.Context) {
	c.JSON(http.StatusOK, authButtons)
}

func registerAuthRoutes(r *gin.Engine) {
	for id, cfg := range authConfig.Backends {
		button := AuthButton{
			DisplayName: cfg.DisplayName,
			Icon: cfg.Icon,
			Color: cfg.Color,
			Path: fmt.Sprintf("/api/oidc/%s/login", id),
		}
		authButtons = append(authButtons, button)
		if cfg.Oidc != nil {
			ctx := context.Background()
			o, err := NewOidcBackend(ctx, id, cfg.Oidc)
			if err != nil {
				log.Fatal(err)
			}
			r.GET(fmt.Sprintf("/api/oidc/%s/login", id), o.Redirect)
			r.GET(fmt.Sprintf("/api/oidc/%s/callback", id), o.Callback)
			continue
		}
	}
	r.GET("/api/auth/buttons", listAuthButtons)
}
