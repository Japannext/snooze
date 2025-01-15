package processor

import (
	"os"
	"fmt"

	"github.com/japannext/snooze/pkg/processor/grouping"
	"github.com/japannext/snooze/pkg/processor/mapping"
	"github.com/japannext/snooze/pkg/processor/notification"
	"github.com/japannext/snooze/pkg/processor/profile"
	"github.com/japannext/snooze/pkg/processor/ratelimit"
	"github.com/japannext/snooze/pkg/processor/silence"
	"github.com/japannext/snooze/pkg/processor/snooze"
	"github.com/japannext/snooze/pkg/processor/transform"
	"github.com/japannext/snooze/pkg/processor/store"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
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

type Config struct {
	Mapping      mapping.Config      `json:",inline" yaml:",inline"`
	Grouping     grouping.Config     `json:",inline" yaml:",inline"`
	Transform    transform.Config    `json:",inline" yaml:",inline"`
	Profile      profile.Config      `json:",inline" yaml:",inline"`
	Silence      silence.Config      `json:",inline" yaml:",inline"`
	Snooze		 snooze.Config       `json:",inline" yaml:",inline"`
	Ratelimit    ratelimit.Config    `json:",inline" yaml:",inline"`
	Notification notification.Config `json:",inline" yaml:",inline"`
	Store		 store.Config        `json:",inline" yaml:",inline"`
}

var (
	env   *Env
)

func initConfig() {
	// Defaults
	viper.SetDefault("PROCESSOR_PIPELINE_FILE", "/etc/snooze/pipeline.yaml")
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

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(env.PipelineFile)
	if err != nil {
		return &Config{}, fmt.Errorf("failed to read %s: %w", env.PipelineFile, err)
	}

	var cfg *Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return &Config{}, fmt.Errorf("failed to interpret yaml in '%s': %w", env.PipelineFile, err)
	}

	return cfg, nil
}
