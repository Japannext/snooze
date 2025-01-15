package alertmanager

import (
	"strconv"
	"strings"
	"time"

	"github.com/japannext/snooze/pkg/common/utils"
	"github.com/japannext/snooze/pkg/models"
	log "github.com/sirupsen/logrus"
)

const (
	identityPrefix = "identity__"
	TIME_LAYOUT     = "2006-01-02T15:04:05.999999999Z"
)

// Convert the alert from Prometheus format to snooze ActiveAlert format.
func convertAlert(alert PostableAlert) *models.ActiveAlert {
	item := &models.ActiveAlert{}

	item.Identity = make(map[string]string)
	item.Labels = make(map[string]string)

	for key, value := range alert.Labels {
		switch {
		case key == "alertname":
			item.AlertName = value
		case key == "alertgroup":
			item.AlertGroup = value
		case key == "severity__text":
			item.SeverityText = value
		case key == "severity_number":
			nb, err := strconv.Atoi(value)
			if err != nil {
				log.Warnf("label `severity__number`: `%s` is not a number", value)

				continue
			}
			item.SeverityNumber = int32(nb)
		case strings.HasPrefix(key, identityPrefix):
			k := strings.TrimPrefix(key, identityPrefix)
			item.Identity[k] = value
		default:
			item.Labels[key] = value
		}
	}

	if item.SeverityNumber == 0 {
		item.SeverityNumber = utils.GuessSeverityNumber(item.SeverityText)
	}

	item.Hash = utils.ComputeHash(alert.Labels)
	item.Source.Kind = ALERTMANAGER_KIND
	item.Source.Name = config.InstanceName
	item.StartsAt = parseTime(alert.StartsAt)

	item.Description = alert.Annotations["description"]
	item.Summary = alert.Annotations["summary"]

	item.LastHit = models.TimeNow()

	return item
}

func parseTime(text string) models.Time {
	t, err := time.Parse(TIME_LAYOUT, text)
	if err != nil {
		log.Warnf("failed to parse time format `%s`: %s", text, err)

		return models.Time{}
	}

	return models.Time{Time: t}
}
