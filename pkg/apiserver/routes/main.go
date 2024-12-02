package routes

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/tracing"
	"github.com/japannext/snooze/pkg/apiserver/auth"
)

var storeQ *mq.Pub
var tracer trace.Tracer

func Register(r *gin.Engine) {
	storeQ = mq.StorePub()
	tracing.Init("snooze-apiserver")
	tracer = otel.Tracer("snooze")

	rr := r.Group("/api")

	// Anonymous
	rr.POST("/alert", postAlert)

	// Authenticated
	authn := rr.Group("/", auth.Authenticate())
	{
		// Ack
		authn.GET("/acks", getAcks)
		authn.POST("/ack", postAck)
		// Alerts
		authn.GET("/alerts", getAlerts)
		// Groups
		authn.GET("/groups", getGroups)
		// Logs
		authn.GET("/logs", getLogs)
		// Notification
		authn.GET("/notifications", getNotifications)
		// Snooze
		authn.GET("/snoozes", getSnoozes)
		authn.POST("/snooze", postSnooze)
		authn.POST("/snooze/cancel", postSnoozeCancel)
		// Tags
		authn.GET("/tags", getTags)
		authn.POST("/tag", postTag)
	}

	// Admin routes
	//adminOnly := rr.Group("/admin", auth.Authenticated(), auth.AdminOnly())
	//{
	//}
}
