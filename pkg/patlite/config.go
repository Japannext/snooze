package patlite

import (
    "github.com/spf13/viper"
    log "github.com/sirupsen/logrus"
)

var env *Env

type Env struct {
    ProfilePath string `mapstructure:"PROFILE_PATH"`
}

func initConfig() {
    viper.SetDefault("PROFILE_PATH", "/config/profiles.yaml")

    viper.AutomaticEnv()
    if err := viper.Unmarshal(&env); err != nil {
        log.Fatal(err)
    }
}
