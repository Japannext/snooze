package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/japannext/snooze/pkg/models"

	// Sub-Processors
	"github.com/japannext/snooze/pkg/processor/grouping"
	"github.com/japannext/snooze/pkg/processor/notification"
	"github.com/japannext/snooze/pkg/processor/ratelimit"
	"github.com/japannext/snooze/pkg/processor/silence"
	"github.com/japannext/snooze/pkg/processor/mapping"
	"github.com/japannext/snooze/pkg/processor/transform"
	"github.com/japannext/snooze/pkg/processor/profile"
)

type Config struct {
	PipelineFile    string `mapstructure:"PROCESSOR_PIPELINE_FILE"`
	DefaultPipeline string `mapstructure:"PROCESSOR_DEFAULT_PIPELINE"`

	ListeningAddress string `mapstructure:"PROCESSOR_LISTENING_ADDRESS"`
	ListeningPort    int    `mapstructure:"PROCESSOR_LISTENING_PORT"`

	// Whether to enable prometheus metrics
	PrometheusEnable bool `mapstructure:"OTEL_PROMETHEUS_ENABLE"`
	// Port the prometheus exporter should listen to
	PrometheusPort int `mapstructure:"OTEL_PROMETHEUS_PORT"`

	// Maximum size of a batch
	BatchSize int `mapstructure:"BATCH_SIZE"`
	// Number of time to wait for messages in batch
	BatchTimeoutSeconds int `mapstructure:"BATCH_TIMEOUT_SECONDS"`

	MaxWorkers int `mapstructure:"MAX_WORKERS"`
}

var config *Config
var pipeline *Pipeline

type Pipeline struct {
	Mappings			[]*mapping.Mapping			 `yaml:"mappings" json:"mappings"`
	Transforms	        []*transform.Transform		 `yaml:"transforms" json:"transforms"`
	Grouping	        []*grouping.Grouping		 `yaml:"groupings" json:"groupings"`
	Profiles			[]*profile.Profile			 `yaml:"profiles" json:"profiles"`
	Silences	        []*silence.Silence			 `yaml:"silences" json:"silences"`
	RateLimits			[]*ratelimit.RateLimit		 `yaml:"ratelimits" json:"ratelimits"`
	Notifications       []*notification.Notification `yaml:"notifications" json:"notifications"`
	DefaultDestinations []models.Destination		 `yaml:"default_destinations" json:"default_destinations"`
}

// A config element that can be loaded and reloaded (when config is live updated)
type Loadable interface {
	LoadConfig() error
}

func initConfig() {
	// Defaults
	viper.SetDefault("PROCESSOR_PIPELINE_FILE", "/etc/snooze/pipeline.yaml")
	viper.SetDefault("PROCESSOR_DEFAULT_PIPELINE", "default")
	viper.SetDefault("PROCESSOR_LISTENING_ADDRESS", "0.0.0.0")
	viper.SetDefault("PROCESSOR_LISTENING_PORT", 8080)
	viper.SetDefault("OTEL_PROMETHEUS_ENABLE", true)
	viper.SetDefault("OTEL_PROMETHEUS_PORT", 9317)
	viper.SetDefault("BATCH_SIZE", 20)
	viper.SetDefault("BATCH_TIMEOUT_SECONDS", 1)
	viper.SetDefault("MAX_WORKERS", 50)

	viper.AutomaticEnv()
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	// Load pipeline
	data, err := os.ReadFile(config.PipelineFile)
	if err != nil {
		log.Fatal(err)
	}
	if err := yaml.Unmarshal(data, &pipeline); err != nil {
		log.Fatal(err)
	}
}
