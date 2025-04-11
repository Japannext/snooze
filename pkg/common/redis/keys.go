package redis

import (
	"fmt"

	"github.com/japannext/snooze/pkg/models"
)

// How ratelimit status keys are named in redis
func RatelimitStatusKey(gr *models.Group) string {
	return fmt.Sprintf("ratelimit_status:%s:%s", gr.Name, gr.Hash)
}

// How snooze keys are named in redis
func SnoozeKey(gr *models.Group) string {
	return fmt.Sprintf("snooze:%s:%s", gr.Name, gr.Hash)
}

// How alertmanager's alerts state are represented in redis
func AlertManagerKey(alertGroup, alertName, hash string) string {
	return fmt.Sprintf("alertmanager:%s:%s:%s", alertGroup, alertName, hash)
}
