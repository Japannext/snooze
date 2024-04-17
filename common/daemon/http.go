package daemon

import (
  "context"
  "net/http"

  "github.com/gin-gonic/gin"
)

type HttpDaemon struct {
  Engine *gin.Engine
  srv *http.Server
}

func NewHttpDaemon() *HttpDaemon {
  router := gin.Default()
  return &HttpDaemon{
    Engine: router,
    srv: &http.Server{
      Addr: "",
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

func (d *HttpDaemon) HandleStop() {
  d.srv.Shutdown(context.Background())
}
