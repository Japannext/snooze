package processor

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Env struct {
	PipelineFile    string `mapstructure:"PROCESSOR_PIPELINE_FILE"`
	DefaultPipeline string `mapstructure:"PROCESSOR_DEFAULT_PIPELINE"`

	ListeningAddress string `mapstructure:"PROCESSOR_LISTENING_ADDRESS"`
	ListeningPort    int    `mapstructure:"PROCESSOR_LISTENING_PORT"`

	// Maximum size of a batch
	BatchSize int `mapstructure:"BATCH_SIZE"`
	// Number of time to wait for messages in batch
	BatchTimeoutSeconds int `mapstructure:"BATCH_TIMEOUT_SECONDS"`

	MaxWorkers int `mapstructure:"MAX_WORKERS"`
}

var (
	env *Env
)

func initConfig() {
	// Defaults
	viper.SetDefault("PROCESSOR_DEFAULT_PIPELINE", "default")
	viper.SetDefault("PROCESSOR_LISTENING_ADDRESS", "0.0.0.0")
	viper.SetDefault("PROCESSOR_LISTENING_PORT", 8080)
	viper.SetDefault("BATCH_SIZE", 20)
	viper.SetDefault("BATCH_TIMEOUT_SECONDS", 1)
	viper.SetDefault("MAX_WORKERS", 50)

	viper.AutomaticEnv()
	if err := viper.Unmarshal(&env); err != nil {
		log.Fatal(err)
	}
}
