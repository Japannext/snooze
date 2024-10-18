package opensearch

import (

	"github.com/spf13/viper"
)

type Config struct {
	Addresses          []string `mapstructure:"OPENSEARCH_ADDRESSES"`
	Username           string   `mapstructure:"OPENSEARCH_USERNAME"`
	Password           string   `mapstructure:"OPENSEARCH_PASSWORD"`
	InsecureSkipVerify bool     `mapstructure:"OPENSEARCH_INSECURE_SKIP_VERIFY"`
}

var config *Config

func initConfig() {
	v := viper.New()

	// Defaults
	v.SetDefault("OPENSEARCH_ADDRESSES", "http://127.0.0.1:9200")
	v.SetDefault("OPENSEARCH_INSECURE_SKIP_VERIFY", false)
	v.BindEnv("OPENSEARCH_USERNAME", "")
	v.BindEnv("OPENSEARCH_PASSWORD", "")

	v.AutomaticEnv()
	if err := v.Unmarshal(&config); err != nil {
		log.Fatalf("failed to load opensearch config: %s", err)
	}
}
