package googlechat

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/japannext/snooze/pkg/common/utils"
)

var config *Config

type Config struct {
	ProfilePath        string          `mapstructure:"GOOGLECHAT_PROFILE_PATH"`
	ServiceAccountPath string          `mapstructure:"GOOGLECHAT_SA_PATH"`
	TLS                utils.TLSConfig `mapstructure:",squash"`
}

func initConfig() {
	viper.SetDefault("GOOGLECHAT_PROFILE_PATH", "/config/profiles.yaml")
	viper.SetDefault("GOOGLECHAT_SA_PATH", "/sa/sa_secrets.json")
	viper.SetDefault("SNOOZE_CACERT", "/cacert/ca.crt")

	viper.AutomaticEnv()
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	if err := config.TLS.Load(); err != nil {
		log.Fatal(err)
	}

}
