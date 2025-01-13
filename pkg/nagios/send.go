package nagios

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/japannext/snooze/pkg/models"
)

func sendAlert(item *models.ActiveAlert) error {
	client := &http.Client{}

	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal alert to json: %w", err)
	}

	alertURL := snoozeURL.JoinPath("/api/alert")
	if resp, err := client.Post(alertURL.String(), "application/json", bytes.NewBuffer(data)); err != nil {
		return fmt.Errorf("error posting alert to '%s': %w\n%s", alertURL, err, resp.Body)
	}
	return nil
}
