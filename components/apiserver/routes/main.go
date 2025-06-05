package routes

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/japannext/snooze/components/apiserver/auth"
	"github.com/japannext/snooze/internal/mq"
)

var (
	storeQ *mq.Pub
	tracer trace.Tracer
)

func Register(r *gin.Engine) {
	storeQ = mq.StorePub()
	tracer = otel.Tracer("snooze")

	rr := r.Group("/api")

	// Anonymous
	rr.POST("/alert", postAlert)

	// Authenticated
	authn := rr.Group("/", auth.Authenticate(), auth.UserAllowed())
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

		authn.GET("/condition-validation", getConditionValidation)

		authn.GET("/pipeline/mappings", getPipelineMappings)
		authn.GET("/pipeline/groupings", getPipelineGroupings)
		authn.GET("/pipeline/transforms", getPipelineTransforms)
		authn.GET("/pipeline/profiles", getPipelineProfiles)
		authn.GET("/pipeline/silences", getPipelineSilences)
		authn.GET("/pipeline/ratelimits", getPipelineRatelimits)
		authn.GET("/pipeline/notifications", getPipelineNotifications)
		authn.GET("/pipeline/default-destinations", getPipelineDefaultDestinations)

		authn.GET("/autocomplete/log", getAutocompleteLog)
	}

	// Admin routes
	// adminOnly := rr.Group("/admin", auth.Authenticated(), auth.AdminOnly())
	//{
	//}
}
