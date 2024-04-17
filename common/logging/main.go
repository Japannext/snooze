package logging

import (
  log "github.com/sirupsen/logrus"
)

func Init() {
  cfg, err := initConfig()
  if err != nil {
    log.Fatal(err)
  }

  ll, err := log.ParseLevel(cfg.Level)
  if err != nil {
    log.Fatalf("Unsupported log level '%s': %s", cfg.Level, err)
  }
  log.SetLevel(ll)
  log.Debug("Log level set to:", ll)
}
