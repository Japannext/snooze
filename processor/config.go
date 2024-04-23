package processor

import (
  "os"

  log "github.com/sirupsen/logrus"
  "gopkg.in/yaml.v3"
  "github.com/spf13/viper"

  // Sub-Processors
  "github.com/japannext/snooze/processor/grouping"
  "github.com/japannext/snooze/processor/notification"
  "github.com/japannext/snooze/processor/ratelimit"
  "github.com/japannext/snooze/processor/silence"
  // "github.com/japannext/snooze/processor/snooze"
  "github.com/japannext/snooze/processor/transform"
)

type Config struct {
  PipelineFile string `mapstructure:"PROCESSOR_PIPELINE_FILE"`
  DefaultPipeline string `mapstructure:"PROCESSOR_DEFAULT_PIPELINE"`

  ListeningAddress string `mapstructure:"PROCESSOR_LISTENING_ADDRESS"`
  ListeningPort int `mapstructure:"PROCESSOR_LISTENING_PORT"`

  // Whether to enable prometheus metrics
  PrometheusEnable bool `mapstructure:"OTEL_PROMETHEUS_ENABLE"`
  // Port the prometheus exporter should listen to
  PrometheusPort int `mapstructure:"OTEL_PROMETHEUS_PORT"`
}

var config *Config
var pipeline *Pipeline

type Pipeline struct {
  Name string `yaml:"name"`
  TransformRules []*transform.Rule `yaml:"transform_rules"`
  GroupingRules []*grouping.Rule `yaml:"grouping_rules"`
  SilenceRules []*silence.Rule `yaml:"silence_rules"`
  RateLimit *ratelimit.Rule `yaml:"ratelimit"`
  NotificationRules []*notification.Rule `yaml:"notification_rules"`
  DefaultNotificationChannels []string `yaml:"default_notification_channels"`
}

func initConfig() {
  // Defaults
  viper.SetDefault("PROCESSOR_PIPELINE_FILE", "/etc/snooze/pipeline.yaml")
  viper.SetDefault("PROCESSOR_DEFAULT_PIPELINE", "default")
  viper.SetDefault("PROCESSOR_LISTENING_ADDRESS", "0.0.0.0")
  viper.SetDefault("PROCESSOR_LISTENING_PORT", 8080)
  viper.SetDefault("OTEL_PROMETHEUS_ENABLE", true)
  viper.SetDefault("OTEL_PROMETHEUS_PORT", 9317)

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
