package apiserver

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/japannext/snooze/pkg/common/opensearch"
    //"github.com/japannext/snooze/pkg/common/tracing"
    "github.com/japannext/snooze/pkg/models"
)

func init() {
    routes = append(routes, func(r *gin.Engine) {
		r.GET("/api/acks", getAcks)
        r.POST("/api/ack", postAck)
    })
}

func getAcks(c *gin.Context) {
}

func postAck(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "postAck")
	defer span.End()

	var item models.Ack
	c.BindJSON(&item)

	req := opensearch.UpdateByQueryReq{
		Index: models.LOG_INDEX,
	}

	req.Doc.WithTerms("logIDs", item.LogIDs)

	_, err := opensearch.UpdateByQuery(ctx, req)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to update in opensearch: %s", err)
		return
	}
}
