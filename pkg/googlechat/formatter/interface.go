package formatter

import (
	chat "google.golang.org/api/chat/v1"

	"github.com/japannext/snooze/pkg/models"
)

type Interface interface {
	Format(*models.Notification) (*chat.Message, error)
}
