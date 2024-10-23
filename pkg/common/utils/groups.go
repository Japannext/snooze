package utils

import (
	"github.com/japannext/snooze/pkg/models"
)

func GetGroup(item *models.Log, name string) (*models.Group, bool) {
	for _, group := range item.Groups {
		if group.Name == name {
			return group, true
		}
	}
	return nil, false
}
