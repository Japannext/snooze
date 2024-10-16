package activecheck

import (
	"github.com/gin-gonic/gin"
)

func probeHandler(c *gin.Context) {
	name := c.Param("name")

	relay, ok := probes[name]
	if !ok {
		// TODO
	}
	relay.ServeHTTP(c)
}
