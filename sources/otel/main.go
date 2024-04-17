package otel

import (
  "context"
  "fmt"
  "os"
  "os/signal"
  "syscall"

  "golang.org/x/sync/errgroup"
  log "github.com/sirupsen/logrus"

  "github.com/japannext/snooze/common/logging"
  "github.com/japannext/snooze/common/opensearch"
)

type Daemon interface {
  Start() error
  Stop()
}

func handleServer(name string, errs *errgroup.Group, ctx context.Context, d Daemon) {
  // Starting the server
  errs.Go(d.Start)
  // Gracefully stopping the server upon context termination
  errs.Go(func() error {
    <-ctx.Done()
    log.Debug(fmt.Sprintf("Stopping %s server...", name))
    d.Stop()
    log.Info(fmt.Sprintf("Stopped %s server", name))
    return nil
  })
}

func handleSignal(ctx context.Context) error {
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
    case <-ctx.Done():
      return nil
  }
}

var db *opensearch.Database
var pq *queue.ProcessQueue

func Run() {

  ctx := context.Background()
  errs, ctx := errgroup.WithContext(ctx)

  pq = queue.InitProcessQueue()

  if err := config.init(); err != nil {
    log.Fatal(err)
  }
  if err := logging.Init(); err != nil {
    log.Fatal(err)
  }
  if err := opensearch.Init(); err != nil {
    log.Fatal(err)
  }
  if err = opensearch.DB.CheckHealth(); err != nil {
    log.Fatal(err)
  }

  // Running daemons
  handleServer("otel-grpc", errs, ctx, NewOtelGrpcServer())

  // Catch SIGTERM signals and exit everything
  errs.Go(func() error {
    return handleSignal(ctx)
  })

  // Waiting for daemons. Will return if one daemon fails
  err = errs.Wait()

  if err == context.Canceled || err == nil {
    log.Info("Gracefully exited server")
  } else if err != nil {
    log.Error(err)
  }

}
