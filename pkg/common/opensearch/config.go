package opensearch

import (
	"crypto/tls"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/spf13/viper"
)

type Config struct {
	Addresses          []string `mapstructure:"OPENSEARCH_ADDRESSES"`
	Username           string   `mapstructure:"OPENSEARCH_USERNAME"`
	Password           string   `mapstructure:"OPENSEARCH_PASSWORD"`
	InsecureSkipVerify bool     `mapstructure:"OPENSEARCH_INSECURE_SKIP_VERIFY"`
}

func initConfig() (opensearch.Config, error) {
	v := viper.New()

	// Defaults
	v.SetDefault("OPENSEARCH_ADDRESSES", "http://127.0.0.1:9200")
	v.SetDefault("OPENSEARCH_INSECURE_SKIP_VERIFY", false)
	v.BindEnv("OPENSEARCH_USERNAME", "")
	v.BindEnv("OPENSEARCH_PASSWORD", "")

	v.AutomaticEnv()
	cfg := &Config{}
	if err := v.Unmarshal(&cfg); err != nil {
		return opensearch.Config{}, err
	}
	config := opensearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.InsecureSkipVerify,
			},
		},
	}
	log.Debugf("OPENSEARCH_ADDRESSES: %s", cfg.Addresses)
	log.Debugf("OPENSEARCH_USERNAME: %s", cfg.Username)
	log.Debugf("OPENSEARCH_PASSWORD: **** (size=%d)", len(cfg.Password))
	return config, nil
}
