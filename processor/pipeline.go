package processor

import (
  "path/filepath"

  "github.com/japannext/snooze/processor/grouping"
  "github.com/japannext/snooze/processor/notification"
  "github.com/japannext/snooze/processor/ratelimit"
  "github.com/japannext/snooze/processor/save"
  "github.com/japannext/snooze/processor/silence"
  "github.com/japannext/snooze/processor/snooze"
  "github.com/japannext/snooze/processor/transform"
)

type Pipeline struct {
  Name string `yaml:"name"`
  TransformRules []transform.Rule `yaml:"transform_rules"`
  GroupingRules []silence.Rule `yaml:"grouping_rules"`
  SilenceRules []silence.Rule `yaml:"silence_rules"`
  RateLimit ratelimit.RateLimit `yaml:"ratelimit"`
  NotificationRules []notification.Rule `yaml:"notification_rules"`
  DefaultNotificationChannels []string `yaml:"default_notification_channels"`
}

func initPipelines() {
  cfg, err := loadPipelineConfig()
  if err != nil {
    log.Fatal(err)
  }
  pipelines, err = buildPipelines(cfg)
  if err != nil {
    log.Fatal(err)
  }
}

func loadPipelineConfig() (map[string]PipelineConfig, error) {
  var cfg map[string]PipelineConfig
  path, err := filepath.Abs(config.PipelinePath)
  if err != nil {
    return cfg, fmt.Errorf("[path %s] %w", path, err)
  }
  var pc map[string]PipelineConfig
  err := filepath.WalkDir(path, func(path string, file fs.DirEntry, err error) error {
    if err != nil {
      return cfg, fmt.Errorf("[path %s] %w", path, err)
    }
    if !file.IsDir() {
      fname := file.Name()
      ext := filepath.Ext(fname)
      if ext == "yml" || ext == "yaml" {
        filepath := filepath.Join(path, fname)
        data, err := os.ReadFile(filepath)
        if err != nil {
          return cfg, fmt.Errorf("[config %s] %w", fname, err)
        }
        var p *PipelineConfig
        if err := yaml.UnMarshal(data, &p); err != nil {
          return cfg, fmt.Errorf("[config %s] %w", fname, err)
        }
        log.Infof("Loaded pipeline '%s'", filepath)
        cfg[p.Name] = p
      }
    }
  })
  if err != nil {
    return cfg, err
  }
  return cfg
}

func buildPipelines(cfg) error {
  for name, pc := range cfg {
    p, err := NewPipeline(pc)
    if err != nil {
      log.Fatalf("[pipeline %s] %w", name, err)
    }
    pipelines[name] = p
  }
}
