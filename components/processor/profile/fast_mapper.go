package profile

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/internal/lang"
	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/pipeline"
	"github.com/japannext/snooze/internal/utils"
)

type FastMapper struct {
	keys   []string
	mapper map[pipeline.Kv][]Profile
	fields map[string]*lang.Field
}

func NewFastMapper(cfgs []pipeline.Profile) (*FastMapper, error) {
	m := map[pipeline.Kv][]Profile{}
	keys := utils.NewOrderedSet[string]()
	fields := map[string]*lang.Field{}

	for _, cfg := range cfgs {
		p, err := NewProfile(cfg)
		if err != nil {
			return &FastMapper{}, fmt.Errorf("error in profile %s: %w", cfg.Name, err)
		}

		m[cfg.Switch] = append(m[cfg.Switch], p)

		value, err := lang.NewField(cfg.Switch.Key)
		if err != nil {
			return &FastMapper{}, fmt.Errorf("invalid field `%s: %w", cfg.Switch.Key, err)
		}

		fields[cfg.Switch.Key] = value

		keys.Append(cfg.Switch.Key)
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

		profiles, found := m.mapper[pipeline.Kv{Key: key, Value: value}]
		if !found {
			continue
		}

		return profiles
	}

	return []Profile{}
}
