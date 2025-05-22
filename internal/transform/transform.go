package transform

import (
	"github.com/japannext/snooze/internal/models"
)

type Transform interface {
	Transform(*models.Log) error
	Load() error
}
