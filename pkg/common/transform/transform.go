package transform

import (
	"github.com/japannext/snooze/pkg/models"
)

type Transform interface {
	Transform(*models.Log) error
	Load() error
}
