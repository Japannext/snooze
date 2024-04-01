package logging

import (
  "fmt"

  log "github.com/sirupsen/logrus"
)

func Init() error {
  cfg, err := initConfig()
  if err != nil {
    return err
  }

  ll, err := log.ParseLevel(cfg.Level)
  if err != nil {
    return fmt.Errorf("Unsupported log level '%s': %w", cfg.Level, err)
  }
  log.SetLevel(ll)
  log.Debug("Log level set to:", ll)
  return nil
}
