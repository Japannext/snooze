package logstore

import (
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

type LogStore interface {
	StoreLog(*api.Log) error
	StoreNotification(*api.Notification) error
	SearchLogs(string, api.TimeRange, api.Pagination) (api.LogsResponse, error)
	GetLog(string) (*api.Alert, error)
}

type BulkStore interface {
	LogStore
	BulkStore([]*api.Log) error
}
