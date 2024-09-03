package v2

type LogStore interface {
	StoreLog(*Log) error
	SearchLogs(string, TimeRange, Pagination) (*LogResults, error)
}

type BulkStore interface {
	LogStore
	BulkStore([]*Log) error
}

type NotificationStore interface {
	StoreNotification(*Notification) error
	SearchNotifications(string, TimeRange, Pagination) (*NotificationResults, error)
}

type AlertStore interface {
	StoreAlert(*Alert) error
	SearchAlerts(string, TimeRange, Pagination) (*AlertResults, error)
}
