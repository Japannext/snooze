package opensearch

import (
  "crypto/tls"
  "net/http"

  "github.com/spf13/viper"
  v2 "github.com/opensearch-project/opensearch-go/v2"
)

type Config struct {
  Addresses []string `mapstructure:"OPENSEARCH_ADDRESSES"`
  Username string `mapstructure:"OPENSEARCH_USERNAME"`
  Password string `mapstructure:"OPENSEARCH_PASSWORD"`
  InsecureSkipVerify bool `mapstructure:"OPENSEARCH_INSECURE_SKIP_VERIFY"`
}

func initConfig() (*v2.Config, error) {
  v := viper.New()

  // Defaults
  v.SetDefault("OPENSEARCH_ADDRESSES", "http://127.0.0.1:9200")
  v.SetDefault("OPENSEARCH_INSECURE_SKIP_VERIFY", false)

  v.AutomaticEnv()
  cfg := &Config{}
  if err := v.Unmarshal(&cfg); err != nil {
    return nil, err
  }
  config := &v2.Config{
    Addresses: cfg.Addresses,
    Username: cfg.Username,
    Password: cfg.Password,
    Transport: &http.Transport{
      TLSClientConfig: &tls.Config{
        InsecureSkipVerify: cfg.InsecureSkipVerify,
      },
    },
  }
  return config, nil
}


