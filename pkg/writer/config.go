package writer

import (
    "github.com/spf13/viper"
    log "github.com/sirupsen/logrus"
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
