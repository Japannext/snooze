package daemon

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type HttpDaemon struct {
	Engine *gin.Engine
	srv    *http.Server
}

type Register = func(*gin.Engine)

func NewHttpDaemon(registers ...Register) *HttpDaemon {
	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/livez", "/readyz", "/metrics"},
	}))
	router.Use(gin.Recovery())
	for _, register := range registers {
		register(router)
	}
	return &HttpDaemon{
		Engine: router,
		srv: &http.Server{
			Addr:    "0.0.0.0:8080",
			Handler: router,
		},
	}
}

func (d *HttpDaemon) Run() error {
	err := d.srv.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (d *HttpDaemon) Stop() {
	log.Debugf("Stopping HTTP server...")
	d.srv.Shutdown(context.Background())
	log.Debug("Stopped HTTP server")
}
