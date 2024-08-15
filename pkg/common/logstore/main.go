package logstore

import (
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

type LogStore interface {
	StoreLog(*api.Log) error
	StoreNotification(*api.Notification) error
	Search(string, api.TimeRange, api.Pagination) ([]*api.Alert, error)
	List(api.Pagination) ([]*api.Alert, error)
	Get(string) (*api.Alert, error)
}

type BulkStore interface {
	LogStore
	BulkStore([]*api.Alert) error
}
