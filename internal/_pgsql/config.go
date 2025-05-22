package pgsql

import (
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

func initConfig() pgx.ConnConfig {
	// https://www.postgresql.org/docs/9.4/libpq-envars.html
	cfg, err := pgx.ParseEnvLibpq()
	if err != nil {
		log.Fatalf("error getting postgresql config: %s", err)
	}
	return cfg
}
