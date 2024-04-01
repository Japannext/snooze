package opensearch

import (
  "crypto/tls"
  "net/http"

  "github.com/spf13/viper"
  v2 "github.com/opensearch-project/opensearch-go/v2"
)

type Config struct {
  Addresses []string `mapstructure:"ADDRESSES"`
  Username string `mapstructure:"USERNAME"`
  Password string `mapstructure:"PASSWORD"`
  InsecureSkipVerify bool `mapstructure:"INSECURE_SKIP_VERIFY"`
}

func (cfg *Config) v2Config() v2.Config {
  return v2.Config{
    Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.InsecureSkipVerify}},
    Addresses: cfg.Addresses,
    Username: cfg.Username,
    Password: cfg.Password,
  }
}

func initConfig() (*Config, error) {
  v := viper.New()
  // Prefix
  v.SetEnvPrefix("OPENSEARCH")

  // Defaults
  v.SetDefault("ADDRESSES", "127.0.0.1:9200")
  v.SetDefault("INSECURE_SKIP_VERIFY", false)

  v.AutomaticEnv()
  cfg := &Config{}
  err := v.Unmarshal(&cfg)
  return cfg, err
}


