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
  Start() error
  Stop()
}

type DaemonRunner struct {
  Errs *errgroup.Group
  Context context.Context
}

func New() *DaemonRunner {
  ctx := context.Background()
  errs, ctx := errgroup.WithContext(ctx)
  return &DaemonRunner{errs, ctx}
}

// An utility to run daemons (process that have a start and a stop)
// gracefully, while handling graceful shutdown.
// Behavior:
// * If one daemon exit, everything should exit gracefully
// * If a SIGTERM is received, everything should exit gracefully
func (dr *DaemonRunner) Main(fn func(context.Context, *DaemonRunner) error) {
  var err error

  // Catch SIGTERM signals and exit everything
  dr.Errs.Go(func() error {
    ch := make(chan os.Signal, 1)
    signal.Notify(ch,
      os.Interrupt,
      syscall.SIGQUIT,
      syscall.SIGTERM,
    )
    select {
      case sig := <-ch:
        log.Errorf("Received signal: %s", sig)
        return fmt.Errorf("Exited due to signal: %s", sig)
      case <-dr.Context.Done():
        return nil
    }
  })

  err = fn(dr.Context, dr)
  if err != nil {
    log.Fatal(err)
  }

  // Waiting for daemons. Will return if one daemon fails
  err = dr.Errs.Wait()

  if err == context.Canceled || err == nil {
    log.Info("Gracefully exited server")
  } else if err != nil {
    log.Error(err)
  }
}

func (dr *DaemonRunner) Run(name string, d Daemon) {
  // Starting the daemon
  dr.Errs.Go(d.Start)
  // Gracefully stopping the server upon context termination
  dr.Errs.Go(func() error {
    <-dr.Context.Done()
    log.Debug(fmt.Sprintf("Stopping %s server...", name))
    d.Stop()
    log.Info(fmt.Sprintf("Stopped %s server", name))
    return nil
  })
}
