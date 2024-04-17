package otel

import (
  "fmt"
  "context"
  "net"

  collectorv1 "go.opentelemetry.io/proto/otlp/collector/logs/v1"
  commonv1 "go.opentelemetry.io/proto/otlp/common/v1"
  healthv1 "google.golang.org/grpc/health/grpc_health_v1"
  log "github.com/sirupsen/logrus"
  logv1 "go.opentelemetry.io/proto/otlp/logs/v1"
  resv1 "go.opentelemetry.io/proto/otlp/resource/v1"

  "github.com/japannext/snooze/common/opensearch"
  "github.com/japannext/snooze/common/api/v2"
)

type server struct {
  collectorv1.UnimplementedLogsServiceServer
}

func (s *server) Export(ctx context.Context, in *collectorv1.ExportLogsServiceRequest) (*collectorv1.ExportLogsServiceResponse, error) {

  var items []v2.AlertEvent
  var failedItems = 0

  rls := in.GetResourceLogs()
  if rls == nil {
    log.Warning("Failed to serve log because it has no resource:", in)
    return &collectorv1.ExportLogsServiceResponse{}, nil
  }
  for _, rl := range rls {
    resource := rl.GetResource()
    sls := rl.GetScopeLogs()
    for _, sl := range sls {
      scope := sl.GetScope()
      lrs := sl.GetLogRecords()
      for _, lr := range lrs {
        item := convertLog(resource, scope, lr)
        if err != nil {
          failedItems++
          continue
        }
        pq.Publish(item)
      }
    }
  }

  
  opensearch.DB.BulkInsertAlertEvents(ctx, items)

  return &collectorv1.ExportLogsServiceResponse{}, nil
}
