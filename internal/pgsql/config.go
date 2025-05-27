package pgsql

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Host     string `mapstructure:"PGHOST"`
	Port     int    `mapstructure:"PGPORT"`
	Database string `mapstructure:"PGDATABASE"`
	Username string `mapstructure:"PGUSER"`
	Password string `mapstructure:"PGPASSWORD"`
}

var cfg *Config
var cfgOnce sync.Once

func GetConfig() *Config {
	cfgOnce.Do(func() {
		v := viper.New()
		v.MustBindEnv("PGHOST")
		v.SetDefault("PGPORT", 5432)
		v.SetDefault("PGDATABASE", "snooze")
		v.MustBindEnv("PGUSER")
		v.MustBindEnv("PGPASSWORD")
		v.AutomaticEnv()
		if err := v.Unmarshal(&cfg); err != nil {
			log.Fatalf("failed to load pgsql config: %s", err)
		}
	})
	return cfg
}
