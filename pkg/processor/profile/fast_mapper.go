package profile

import (
	"fmt"

	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/common/utils"
	"github.com/japannext/snooze/pkg/models"
	log "github.com/sirupsen/logrus"
)

type Kv struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type FastMapper struct {
	keys   []string
	mapper map[Kv][]Profile
	fields map[string]*lang.Field
}

func NewFastMapper(profiles []Profile) (*FastMapper, error) {
	m := map[Kv][]Profile{}
	keys := utils.NewOrderedSet[string]()
	fields := map[string]*lang.Field{}

	for _, prf := range profiles {
		m[prf.Switch] = append(m[prf.Switch], prf)

		value, err := lang.NewField(prf.Switch.Key)
		if err != nil {
			return &FastMapper{}, fmt.Errorf("invalid field `%s: %w", prf.Switch.Key, err)
		}

		fields[prf.Switch.Key] = value

		keys.Append(prf.Switch.Key)
	}

	return &FastMapper{keys.Items(), m, fields}, nil
}

func (m *FastMapper) GetMatches(item *models.Log) []Profile {
	for _, key := range m.keys {
		field, found := m.fields[key]
		if !found {
			log.Warnf("unexpected field `%s`", key)

			continue
		}

		value, err := lang.ExtractField(item, field)
		if err != nil {
			continue
		}

		profiles, found := m.mapper[Kv{key, value}]
		if !found {
			continue
		}

		return profiles
	}

	return []Profile{}
}
