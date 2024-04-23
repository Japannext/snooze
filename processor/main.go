package processor

import (
	"github.com/japannext/snooze/common/daemon"
	"github.com/japannext/snooze/common/health"
	"github.com/japannext/snooze/common/logging"
	"github.com/japannext/snooze/common/opensearch"
	"github.com/japannext/snooze/common/rabbitmq"
	"github.com/japannext/snooze/common/redis"

	// Sub-Processors
	"github.com/japannext/snooze/processor/grouping"
	"github.com/japannext/snooze/processor/notification"
	"github.com/japannext/snooze/processor/ratelimit"
	"github.com/japannext/snooze/processor/silence"
	"github.com/japannext/snooze/processor/snooze"
	"github.com/japannext/snooze/processor/transform"
)

var ch *rabbitmq.ProcessChannel

func Run() {

	logging.Init()
	initConfig()
	opensearch.Init()
	redis.Init()
	rabbitmq.Init()
	ch = rabbitmq.InitProcessChannel()

	transform.InitRules(pipeline.TransformRules)
	silence.InitRules(pipeline.SilenceRules)
	grouping.InitRules(pipeline.GroupingRules)
	snooze.Init()
	ratelimit.Init(pipeline.RateLimit)
	notification.InitRules(pipeline.NotificationRules, pipeline.DefaultNotificationChannels)

	dm := daemon.NewDaemonManager()
	p := &Processor{}
	dm.Add("processor", p)
	h := health.HealthStatus{}
	srv := daemon.NewHttpDaemon()
	srv.Engine.GET("/livez", h.LiveRoute)
	srv.Engine.GET("/readyz", h.ReadyRoute)
	dm.Add("health", srv)

	h.Live()
	h.Ready()
	dm.Run()
}
