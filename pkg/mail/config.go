package mail

import (
	"github.com/japannext/snooze/pkg/common/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var config *Config

type Config struct {
	Server        string          `mapstructure:"SMTP_SERVER"`
	Port          int             `mapstructure:"SMTP_PORT"`
	Queue         string          `mapstructure:"SMTP_QUEUE"`
	DefaultSender string          `mapstructure:"SMTP_DEFAULT_SENDER"`
	ProfilePath   string          `mapstructure:"SMTP_PROFILE_PATH"`
	TLS           utils.TLSConfig `mapstructure:",squash"`
}

func initConfig() {
	viper.BindEnv("SMTP_SERVER")
	viper.BindEnv("SMTP_DEFAULT_SENDER")
	viper.SetDefault("SMTP_QUEUE", "mail")
	viper.SetDefault("SMTP_PORT", "25")
	viper.SetDefault("SMTP_PROFILE_PATH", "/config/profiles.yaml")
	viper.SetDefault("SNOOZE_CACERT", "/cacert/ca.crt")

	viper.AutomaticEnv()
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	if err := config.TLS.Load(); err != nil {
		log.Fatal(err)
	}
}
