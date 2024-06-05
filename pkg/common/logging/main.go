package logging

import (
	log "github.com/sirupsen/logrus"
)

func Init() {
	cfg := initConfig()

	ll, err := log.ParseLevel(cfg.Level)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(ll)
	log.Debug("Log level set to:", ll)
}
