package formatter

import (
	"github.com/japannext/snooze/internal/models"
	chat "google.golang.org/api/chat/v1"
)

type Interface interface {
	Format(*models.Notification) (*chat.Message, error)
}
