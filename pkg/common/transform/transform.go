package transform

import (
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

type Transform interface {
	Transform(*api.Log) error
	Load() error
}
