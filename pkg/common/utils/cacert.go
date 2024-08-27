package utils

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

type TLSConfig struct {
	CaCert string `mapstructure:"SNOOZE_CACERT" yaml:"cacert"`

	internal struct {
		config *tls.Config
	}
}

func (cfg *TLSConfig) Load() error {
	data, err := ioutil.ReadFile(cfg.CaCert)
	if err != nil {
		return err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(data)
	cfg.internal.config = &tls.Config{RootCAs: pool}
	return nil
}

func (cfg *TLSConfig) Config() *tls.Config {
	return cfg.internal.config
}
