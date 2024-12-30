package activecheck

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/japannext/snooze/pkg/models"
)

// When the notification encountered is an active check,
// return the callback for the active check to count in.
func Callback(item *models.Notification) error {
	if item.ActiveCheckURL == "" {
		return nil
	}
	callback := &models.SourceActiveCheck{}
	data, err := json.Marshal(callback)
	if err != nil {
		return err
	}
	if _, err := http.Post(item.ActiveCheckURL, "application/json", bytes.NewBuffer(data)); err != nil {
		return err
	}
	return nil
}
