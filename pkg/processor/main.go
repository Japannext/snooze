package processor

import (
	"github.com/japannext/snooze/pkg/common/daemon"
	"github.com/japannext/snooze/pkg/common/health"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/rabbitmq"
	"github.com/japannext/snooze/pkg/common/redis"

	// Sub-Processors
	"github.com/japannext/snooze/pkg/processor/grouping"
	"github.com/japannext/snooze/pkg/processor/notification"
	"github.com/japannext/snooze/pkg/processor/ratelimit"
	"github.com/japannext/snooze/pkg/processor/silence"
	"github.com/japannext/snooze/pkg/processor/snooze"
	"github.com/japannext/snooze/pkg/processor/transform"
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
	dm.Add("http", srv)

	h.Live()
	h.Ready()
	dm.Run()
}
