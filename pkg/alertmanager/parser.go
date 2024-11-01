package alertmanager

import (
	"strings"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/common/utils"
	"github.com/japannext/snooze/pkg/models"
)

const IDENTITY_PREFIX = "identity__"
const TIME_LAYOUT = "2006-01-02T15:04:05.999999999Z"

func parseIdentity(labels map[string]string) map[string]string {
	identity := map[string]string{}
	for key, value := range labels {
		if strings.HasPrefix(key, IDENTITY_PREFIX) {
			k := strings.TrimPrefix(key, IDENTITY_PREFIX)
			identity[k] = value
		}
	}
	return identity
}

func parseSeverity(labels map[string]string) (string, int32) {
	text, ok := labels["severity__text"]
	if !ok {
		text = labels["severity"]
	}

	var number int32
	if n, ok := labels["severity__number"]; ok {
		i, err := strconv.Atoi(n)
		if err != nil {
			log.Warnf("in label `severity__number`, `%s` is not a number", n)
		}
		number = int32(i)
	}

	if number == 0 {
		number = utils.GuessSeverityNumber(text)
	}

	return text, number
}

func parseTime(text string) models.Time {
	t, err := time.Parse(TIME_LAYOUT, text)
	if err != nil {
		log.Warnf("failed to parse time format `%s`: %s", text, err)
		return models.Time{}
	}
	return models.Time{Time: t}
}
