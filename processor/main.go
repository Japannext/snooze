package processor

import (
  "github.com/japannext/snooze/common/daemon"
  "github.com/japannext/snooze/common/opensearch"
  "github.com/japannext/snooze/common/redis"
  "github.com/japannext/snooze/common/health"
  "github.com/japannext/snooze/common/rabbitmq"

  // Sub-Processors
  "github.com/japannext/snooze/processor/grouping"
  "github.com/japannext/snooze/processor/notification"
  "github.com/japannext/snooze/processor/ratelimit"
  "github.com/japannext/snooze/processor/silence"
  "github.com/japannext/snooze/processor/snooze"
  "github.com/japannext/snooze/processor/transform"
)

func Run() {

  logging.Init()
  initConfig()
  opensearch.Init()
  redis.Init()
  rabbitmq.Init()

  transform.InitRules(config.TransformRules)
  silence.InitRules(config.SilenceRules)
  grouping.InitRules(config.GroupingRules)
  snooze.Init()
  ratelimit.Init(config.RateLimit)
  notification.InitRules(config.NotificationRules, config.DefaultNotificationChannels)

  dm := daemon.NewDaemonManager()
  dm.Add("consumer", newConsumer())
  h := health.Health{}
  srv := daemon.NewHttpDaemon()
  srv.Engine.GET("/livez", h.LiveRoute())
  srv.Engine.GET("/readyz", h.ReadyRoute())
  dm.Add("health", srv)

  h.Live()
  h.Ready()
  dm.Run()
}
