package routes

import (
	// "net/http".
	"time"

	"github.com/gin-gonic/gin"
)

type Timeline struct {
	Events []Event
}

type Event struct {
	Time        time.Time
	Label       string
	Description string
	Color       string
}

func logTimeline(c *gin.Context) {
}

/*
func alertTimeline(c *gin.Context) {

	uid := c.Param("uid")

	req := opensearchapi.MSearchReq{
		Indices: []string{
			api.ALERT_INDEX,
			api.NOTIFICATION_INDEX,
		},
		Body: body,
	}
	resp, err := opensearch.Client.MSearch(c, req)
	if err != nil {
		// TODO
		return
	}
	resp.Hits

}
*/
