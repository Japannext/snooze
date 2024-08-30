package daemon

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"golang.org/x/sync/errgroup"
)

type Daemon interface {
	Run() error
	Stop()
}

type DaemonManager struct {
	errs    *errgroup.Group
	context context.Context
	daemons map[string]Daemon
	ready []Check
	live []Check
}

func NewDaemonManager() *DaemonManager {
	ctx := context.Background()
	errs, ctx := errgroup.WithContext(ctx)
	return &DaemonManager{
		errs: errs,
		context: ctx,
		daemons: map[string]Daemon{},
		ready: []Check{},
		live: []Check{},
	}
}

func (dm *DaemonManager) AddDaemon(name string, d Daemon) {
	dm.daemons[name] = d
}

func (dm *DaemonManager) Run() {
	var err error

	dm.setupHealthcheck()
	dm.setupPrometheus()

	// Catch SIGTERM signals and exit everything
	dm.errs.Go(func() error {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit,
			os.Interrupt,
			syscall.SIGQUIT,
			syscall.SIGTERM,
		)
		select {
		case sig := <-exit:
			log.Errorf("Received signal: %s", sig)
			return fmt.Errorf("Exited due to signal: %s", sig)
		case <-dm.context.Done():
			return nil
		}
	})

	for name, _ := range dm.daemons {
		n := name
		d := dm.daemons[n]
		log.Debug(fmt.Sprintf("Starting '%s' daemon", n))
		dm.errs.Go(d.Run)
		dm.errs.Go(func() error {
			<-dm.context.Done()
			log.Debug(fmt.Sprintf("Stopping '%s' daemon...", n))
			d.Stop()
			log.Info(fmt.Sprintf("Stopped '%s' daemon", n))
			return nil
		})
	}

	// Waiting for daemons. Will return if one daemon fails
	err = dm.errs.Wait()
	if err == context.Canceled || err == nil {
		log.Info("Gracefully exited server")
	} else if err != nil {
		log.Error(err)
	}
}
