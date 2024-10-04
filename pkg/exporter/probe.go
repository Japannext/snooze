package exporter

import (
	"time"
	"sync"

	"github.com/gin-gonic/gin"
)

func waitWithTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()
	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

func probeHandler(c *gin.Context) {
	name := c.Param("name")

	relay, ok := probes[name]
	if !ok {
		// TODO
	}
	relay.ServeHTTP(c)
}
