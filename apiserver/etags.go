package apiserver

import (
  "crypto/md5"
  "fmt"
  "io"
  "os"
  "path/filepath"

  log "github.com/sirupsen/logrus"
  "github.com/gin-gonic/gin"
)

type eTagCache struct {
  paths map[string]string
}

func (e *eTagCache) build() {
  filepath.Walk(config.StaticPath, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      log.Warnf("Could not generate eTag cache for %s: %s", path, err)
    }
    f, err := os.Open(path)
    if err != nil {
      log.Warnf("Could not generate eTag cache for %s: %s", path, err)
    }
    defer f.Close()
    h := md5.New()
    if _, err := io.Copy(h, f); err != nil {
      log.Warnf("Could not generate eTag cache for %s: %s", path, err)
    }
    e.paths[path] = fmt.Sprintf("%x", h.Sum(nil))
    return nil
  })
}

func (e *eTagCache) get(p string) (string, bool) {
  value, ok := e.paths[p]
  return value, ok
}

func eTagMiddleware() gin.HandlerFunc {
  etags := &eTagCache{}
  etags.build()
  return func(c *gin.Context) {
    c.Header("Cache-Control", "public max-age=31436000")
    path := c.Param("filepath") // from the .Static() route
    etag, ok := etags.get(path)
    if ok {
      c.Header("ETag", etag)
    }
    c.Next()
  }
}

