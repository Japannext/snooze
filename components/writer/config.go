package writer

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	BatchSize int `mapstructure:"BATCH_SIZE"`
}

var config *Config

func initConfig() {
	viper.SetDefault("BATCH_SIZE", 100)

	viper.AutomaticEnv()
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
}
