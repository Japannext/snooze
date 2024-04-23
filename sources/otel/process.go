package otel

import (
	"context"

	log "github.com/sirupsen/logrus"
	collectorv1 "go.opentelemetry.io/proto/otlp/collector/logs/v1"
)

type server struct {
	collectorv1.UnimplementedLogsServiceServer
}

func (s *server) Export(ctx context.Context, in *collectorv1.ExportLogsServiceRequest) (*collectorv1.ExportLogsServiceResponse, error) {

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
				alert := convertLog(resource, scope, lr)
				if err != nil {
					failedItems++
					continue
				}
				pq.Publish(ctx, alert)
			}
		}
	}

	return &collectorv1.ExportLogsServiceResponse{}, nil
}
