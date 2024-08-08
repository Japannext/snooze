package logstore

import (
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

type LogStore interface {
	Store(*api.Alert) error
	Search(string, api.Pagination) ([]*api.Alert, error)
	List(api.Pagination) ([]*api.Alert, error)
}

type BulkStore interface {
	LogStore
	BulkStore([]*api.Alert) error
}
