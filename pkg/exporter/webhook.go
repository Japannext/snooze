package exporter

import (
	"github.com/gin-gonic/gin"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func webhookHandler(c *gin.Context) {
	key := c.Param("uid")
	var callback *api.ActiveCheckCallback
	c.BindJSON(&callback)

	waiter.Insert(key, callback)
}
