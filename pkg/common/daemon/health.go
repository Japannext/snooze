package daemon

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Check interface {
	Name() string
	Pass() error
}

func (dm *DaemonManager) getGinEngine() *gin.Engine {
	d, found := dm.daemons["http"]
	if !found {
		d = NewHttpDaemon()
		dm.daemons["http"] = d
	}
	dd, ok := d.(*HttpDaemon)
	if !ok {
		log.Fatal("cannot define a 'http' daemon without it being a HttpDaemon")
	}
	return dd.Engine
}

func (dm *DaemonManager) setupHealthcheck() {
	r := dm.getGinEngine()
	if len(dm.ready) == 0 {
		dm.ready = append(dm.ready, NewPing())
	}
	if len(dm.live) == 0 {
		dm.live = append(dm.live, NewPing())
	}
	r.GET("/livez", dm.livezRoute)
	r.GET("/readyz", dm.readyzRoute)
}

func (dm *DaemonManager) livezRoute(c *gin.Context) {
	var (
		status = http.StatusOK
		err error
	)
	var buf strings.Builder
	for _, check := range dm.live {
		err = check.Pass()
		if err != nil {
			status = http.StatusInternalServerError
			buf.WriteString("[X] ")
			buf.WriteString(check.Name())
			buf.WriteString(": ")
			buf.WriteString(err.Error())
			buf.WriteString("\n")
			continue
		}
		buf.WriteString("[✅] ")
		buf.WriteString(check.Name())
		buf.WriteString("\n")
	}
	c.String(status, buf.String())
}

func (dm *DaemonManager) readyzRoute(c *gin.Context) {
	var (
		status = http.StatusOK
		err error
	)
	var buf strings.Builder
	for _, check := range dm.ready {
		err = check.Pass()
		if err != nil {
			status = http.StatusInternalServerError
			buf.WriteString("[X] ")
			buf.WriteString(check.Name())
			buf.WriteString(": ")
			buf.WriteString(err.Error())
			buf.WriteString("\n")
			continue
		}
		buf.WriteString("[✅] ")
		buf.WriteString(check.Name())
		buf.WriteString("\n")
	}
	c.String(status, buf.String())
}

func (dm *DaemonManager) AddLiveCheck(check Check) {
	dm.live = append(dm.live, check)
}

func (dm *DaemonManager) AddReadyCheck(check Check) {
	dm.live = append(dm.ready, check)
}

