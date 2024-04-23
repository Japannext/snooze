package otel

import (
  "fmt"
  "net"

  log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

  collectorv1 "go.opentelemetry.io/proto/otlp/collector/logs/v1"
  "google.golang.org/grpc/health"
  healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

var err error

type OtelGrpcServer struct {
  grpc *grpc.Server
}

func NewOtelGrpcServer() *OtelGrpcServer {
  return &OtelGrpcServer{
    grpc: grpc.NewServer(),
  }
}

func (s *OtelGrpcServer) Start() error {
  log.Infof("Starting opentelemetry listener")

  hostport := fmt.Sprintf("%s:%d", config.GrpcListeningAddress, config.GrpcListeningPort)
  lis, err := net.Listen("tcp", hostport)
  if err != nil {
    return fmt.Errorf("Failed to listen to '%s': %w", hostport, err)
  }

  collectorv1.RegisterLogsServiceServer(s.grpc, &server{})
  healthv1.RegisterHealthServer(s.grpc, health.NewServer())
  if err = s.grpc.Serve(lis); err != nil {
    return fmt.Errorf("Error while serving on %s: %w", hostport, err)
  }

  return nil
}

func (s *OtelGrpcServer) Stop() {
  s.grpc.GracefulStop()
}

