package server

import (
  "context"
  "fmt"
  "os"
  "os/signal"
  "syscall"

  "golang.org/x/sync/errgroup"
  log "github.com/sirupsen/logrus"
)

func initLogging() error {
  ll, err := log.ParseLevel(Config.LogLevel)
  if err != nil {
    return fmt.Errorf("Unsupported log level '%s': %w", Config.LogLevel, err)
  }
  log.SetLevel(ll)
  log.Debug("Log level set to:", ll)
  return nil
}

type Server interface {
  Start() error
  Stop()
}

func handleServer(name string, errs *errgroup.Group, ctx context.Context, s Server) {
  // Starting the server
  errs.Go(s.Start)
  // Gracefully stopping the server upon context termination
  errs.Go(func() error {
    <-ctx.Done()
    log.Debug(fmt.Sprintf("Stopping %s server...", name))
    s.Stop()
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

func Run() {

  ctx := context.Background()
  errs, ctx := errgroup.WithContext(ctx)

  if err := initConfig(); err != nil {
    log.Fatal(err)
  }
  if err := initLogging(); err != nil {
    log.Fatal(err)
  }

  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)

  // Running daemons
  handleServer("otel-grpc", errs, ctx, NewOtelGrpcServer())
  errs.Go(func() error {
    return handleSignal(ctx)
  })

  // Waiting for daemons. Will return if one daemon fails
  err := errs.Wait()

  if err == context.Canceled || err == nil {
    log.Info("Gracefully exited server")
  } else if err != nil {
    log.Error(err)
  }

}
