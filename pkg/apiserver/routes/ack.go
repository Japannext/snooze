package routes

import (
    "net/http"

    "github.com/gin-gonic/gin"
	//log "github.com/sirupsen/logrus"

    "github.com/japannext/snooze/pkg/common/opensearch"
    "github.com/japannext/snooze/pkg/models"
)

func init() {
    registers = append(registers, func(r *gin.Engine) {
		r.GET("/api/acks", getAcks)
        r.POST("/api/ack", postAck)
    })
}

func getAcks(c *gin.Context) {
}

func postAck(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "postAck")
	defer span.End()

	var item *models.Ack
	c.BindJSON(&item)

	item.Time = models.TimeNow()

	if len(item.LogIDs) == 0 {
		c.String(http.StatusBadRequest, "ack object need a non-empty logIDs")
		return
	}

	req := opensearch.UpdateByQueryReq{
		Index: models.LOG_INDEX,
	}
	req.Doc.WithTerms("_id", item.LogIDs)
	req.Doc.WithPainlessScript("ctx._source.status.kind = params.kind", map[string]interface{}{
		"kind": "acked",
	})

	err := opensearch.UpdateByQuery(ctx, req)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to update in opensearch: %s", err)
		return
	}

	err = opensearch.Index(ctx, &opensearch.IndexReq{
		Index: models.ACK_INDEX,
		Item: item,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to save ack to opensearch: %s", err)
		return
	}
}
