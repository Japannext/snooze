package exporter

import (
	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/pkg/models"
)

func webhookHandler(c *gin.Context) {
	key := c.Param("uid")
	var callback *models.ActiveCheckCallback
	c.BindJSON(&callback)

	waiter.Insert(key, callback)
}
