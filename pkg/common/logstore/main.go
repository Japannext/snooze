package logstore

import (
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

type LogStore interface {
	GetLog(string) (*api.Log, error)
	StoreLog(*api.Log) error
	StoreNotification(*api.Notification) error
	SearchLogs(string, api.TimeRange, api.Pagination) (*api.LogResults, error)
}

type BulkStore interface {
	LogStore
	BulkStore([]*api.Log) error
}
